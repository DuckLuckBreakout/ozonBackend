import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import Bus from '../utils/bus/bus';
import HTTPResponses from '../utils/http-responses/httpResponses';

/**
 * @description Model for Product in MVP Arch
 */
class CartModel extends BaseModel {
    #ids
    #products
    #price
    #needsRerender;
    #lastAddedProductID

    /**
     * @return {Object} item
     */
    get products() {
        return this.#products;
    }

    /**
     * @return {Object} price
     */
    get price() {
        return this.#price;
    }

    /**
     * @return {number} ID of last added product (for authentication and future adding)
     */
    get lastAddedProduct() {
        return this.#lastAddedProductID;
    }

    /**
     * @param {number|undefined} newProduct
     */
    set lastAddedProduct(newProduct) {
        this.#lastAddedProductID = newProduct;
    }

    /**
     * @return {Set}
     */
    get ids() {
        return this.#ids;
    }

    /**
     * @return {boolean} Is View needs to be rerendered.
     */
    get needsRerender() {
        return this.#needsRerender;
    }

    /**
     * @description drops the cart
     */
    dropCart = () => {
        this.#ids = undefined;
        this.#products = undefined;
        this.#lastAddedProductID = undefined;
    }

    /**
     *
     * @param {number} id id of product
     * @param {number | string} count amount of product
     * @param {Events} callbackEvent
     */
    addProduct(id, count, callbackEvent) {
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.cartProduct,
            body: {product_id: +id,
                count: +count,
            },
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then(() => {
            Bus.globalBus.emit(Events.HeaderChangeCartItems, count);
            Bus.globalBus.emit(callbackEvent, id);
            (this.#products = this.products || []).push({
                count: 1,
                id: id,
            });
            this.#ids.add(+id);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.#lastAddedProductID = id;
                this.bus.emit(Events.CartProductAdded, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.CartProductAdded, Responses.Offline);
                break;
            }
            default: {
                console.error(err);
                break;
            }
            }
        });
    }

    /**
     *
     * @param {string} event
     */
    getIDs = (event) => {
        if (this.#ids) {
            Bus.globalBus.emit(event, this.#ids);
            return;
        }

        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.cart,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            this.#ids = new Set((parsedJson.products || []).map((elem) => {
                return elem.id;
            }) || []);
            Bus.globalBus.emit(event, this.#ids);
        }).catch((err) => {
            console.error(err);
            Bus.globalBus.emit(event, new Set());
        });
    }

    /**
     * @param {number | string} itemID
     * @param {number | string} count
     */
    changeCartCounter = (itemID, count) => {
        for (const product of this.#products) {
            if (+product.id === +itemID) {
                const diff = +count - this.#products[this.#products.indexOf(product)].count;
                this.#products[this.#products.indexOf(product)].count = +count;
                Bus.globalBus.emit(Events.HeaderChangeCartItems, diff);
                break;
            }
        }
    }

    /**
     * @param {number} id id of product
     * @param {number | string} count new amount of product
     */
    changeItemAmount(id, count) {
        if (count === 0) {
            this.removeProduct(id);
            return;
        }

        this.changeCartCounter(id, count);

        AjaxModule.putUsingFetch({
            url: serverApiPath + urls.cartProduct,
            body: {product_id: +id,
                count: +count},
        }).then((response) => {
            return response.json();
        }).then(() => {
            this.#needsRerender = true;
        }).catch((err) => {
            console.error(err);
        });
    }

    /**
     * @param {number} id id of removed product
     * @param {Events} callbackEvent
     */
    removeProduct = (id, callbackEvent = Events.DummyEvent) => {
        if (!(this.#ids.has(+id))) {
            return;
        }
        this.changeCartCounter(id, 0);
        this.#ids.delete(+id);
        Bus.globalBus.emit(Events.CartProductRemoved, Responses.Success, id);
        Bus.globalBus.emit(callbackEvent, id);
        this.products.splice(this.products.findIndex((elem) => +elem.id === +id), 1);

        AjaxModule.deleteUsingFetch({
            url: serverApiPath + urls.cartProduct,
            body: {product_id: +id},
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Offline: {
                this.bus.emit(Events.CartProductRemoved, Responses.Offline);
                break;
            }
            default: {
                Bus.globalBus.emit(Events.CartProductRemoved, Responses.Error);
            }
            }
        });
    }

    /**
     * @description loads product from cart
     */
    loadProducts() {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.cart,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            this.#needsRerender = false;
            this.#products = parsedJson.products || [];
            this.#price = parsedJson.price;
            this.#ids = new Set(this.#products.map((elem) => {
                return elem.id;
            }) || []);
            for (const product of this.#products) {
                product.preview_image = product.preview_image;
            }
            this.bus.emit(Events.CartLoaded, Responses.Success);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.CartLoaded, Responses.Unauthorized);
                return;
            }
            case HTTPResponses.InternalError: {
                this.#needsRerender = false;
                this.#products = [];
                this.#ids = new Set();
                this.bus.emit(Events.CartLoaded, Responses.Success);
                return;
            }
            }
            this.bus.emit(Events.CartLoaded, Responses.Error);
        });
    }

    loadProductsAmount = () => {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.cart,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            this.#products = parsedJson.products || [];
            this.#ids = new Set(this.#products.map((elem) => {
                return elem.id;
            }) || []);
            let count = 0;
            for (const product of this.#products) {
                count += product.count;
            }
            this.bus.emit(Events.CartLoadedProductsAmountReaction, count);
        }).catch((err) => {
            this.bus.emit(Events.CartLoadedProductsAmountReaction, 0);
        });
    }

    /**
     * @param {number|string} productID
     * @return {boolean}
     */
    isCartContains = (productID) => {
        return this.products.find((elem) => {
            return +elem.id === +productID;
        }) !== undefined;
    }
}

export default CartModel;
