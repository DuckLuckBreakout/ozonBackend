import BasePresenter from './BasePresenter.js';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import Bus from '../utils/bus/bus';
import Router from '../utils/router/Router';

/**
 * @description Presenter for Login View and Model
 */
class CartPresenter extends BasePresenter {
    /**
     * @param {HTMLElement} application html of application
     * @param {Class} View Class of view object
     * @param {Class} Model Class of model object
     */
    constructor(application, View, Model) {
        super(application, View, Model);
        this.bus.on(Events.CartLoad, this.loadProducts);
        this.bus.on(Events.CartLoaded, this.cartLoadedReaction);
        this.bus.on(Events.CartProductAdded, this.cartProductAddedReaction);
        this.bus.on(Events.CartLoadedProductsAmountReaction, this.productsAmountLoadedReaction);

        Bus.globalBus.on(Events.CartProductRemoved, this.productRemovedReaction);
        Bus.globalBus.on(Events.CartAddProduct, this.addProduct);
        Bus.globalBus.on(Events.CartRemoveProduct, this.removeProduct);
        Bus.globalBus.on(Events.CartProductChange, this.changeProduct);
        Bus.globalBus.on(Events.CartLoadProductsAmount, this.loadProductsAmount);
        Bus.globalBus.on(Events.CartGetProductsID, this.getProductsIDs);
        Bus.globalBus.on(Events.CartGetProductID, this.getProductID);
        Bus.globalBus.on(Events.CartAddLastProduct, this.addLastProduct);
        Bus.globalBus.on(Events.OrderSent, this.dropCart);
        Bus.globalBus.on(Events.CartDrop, this.dropCart);
    }

    /**
     * @return {Object} array of products
     */
    get products() {
        return this.model.products;
    }

    /**
     * @return {Object} price of users cart
     */
    get price() {
        return this.model.price;
    }

    /**
     * @return {boolean} is view needs to be rerendered
     */
    get needsRerender() {
        return this.model.needsRerender;
    }

    /**
     * @description loads all products
     */
    loadProducts = () => {
        this.model.loadProducts();
    }

    /**
     * @param {Responses} result
     */
    cartLoadedReaction = (result) => {
        switch (result) {
        case Responses.Success: {
            this.view.render();
            return;
        }
        case Responses.Offline: {
            Router.open('/offline', {replaceState: true});
            break;
        }
        case Responses.Unauthorized: {
            Router.open('/login', {replaceState: true});
            break;
        }
        default: {
            console.error(result);
            break;
        }
        }
    }

    /**
     *
     * @param {number} id
     * @param {number|string} count
     * @param {Events} callbackEvent
     */
    addProduct = (id, count, callbackEvent) => {
        if (count < 0 || id < 0) {
            return;
        }
        this.model.addProduct(id, count, callbackEvent);
    }

    changeProduct = ({id, count}) => {
        this.model.changeItemAmount(id, count);
    }

    /**
     * @param {number} id
     * @param {Events} callbackEvent
     */
    removeProduct = (id, callbackEvent) => {
        this.model.removeProduct(id, callbackEvent);
    }

    /**
     * @param {string} result
     * @param {number} id
     */
    productRemovedReaction = (result, id = undefined) => {
        if (result === Responses.Success) {
            this.view.removeProduct(id);
            return;
        }
        console.error(result);
    }

    /**
     * @param {string} result
     */
    cartProductAddedReaction = (result) => {
        switch (result) {
        case Responses.Success: {
            break;
        }
        case Responses.Offline: {
            Router.open('/offline');
            break;
        }
        case Responses.Unauthorized: {
            Router.open('/login');
            break;
        }
        default: {
            console.error(result);
            break;
        }
        }
    }

    loadProductsAmount = () => {
        this.model.loadProductsAmount();
    }

    /**
     * @param {number} count
     */
    productsAmountLoadedReaction = (count) => {
        Bus.globalBus.emit(Events.HeaderSetCartItems, count);
        this.addLastProduct();
    }

    /**
     * @param {string} eventToEmit
     */
    getProductsIDs = (eventToEmit) => {
        this.model.getIDs(eventToEmit);
    }

    getProductID = () => {
        this.model.getIDs(Events.CartLoadedProductID);
    }

    addLastProduct = () => {
        if (this.model.lastAddedProduct === undefined || this.model.isCartContains(this.model.lastAddedProduct)) {
            return;
        }
        this.model.addProduct(this.model.lastAddedProduct, 1, Events.ProductsItemAdded);
        this.model.addProduct(this.model.lastAddedProduct, 0, Events.ProductItemAdded);
        this.model.lastAddedProduct = undefined;
    }

    dropCart = () => {
        this.model.dropCart();
    }
}

export default CartPresenter;
