import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';
import Router from '../utils/router/Router';

/**
 * @description Model for Product in MVP Arch
 */
class OrderModel extends BaseModel {
    #products
    #price
    #recipient
    #address
    #promo

    /**
     *
     * @return {Array} array of products
     */
    get products() {
        return this.#products;
    }

    /**
     * @return {Object} total price, discount price, discount size
     */
    get price() {
        return this.#price;
    }

    /**
      * @return {Object} Information about recipient
     */
    get recipient() {
        return this.#recipient;
    }

    /**
     * @return {string} Address of delivery
     */
    get address() {
        return this.#address;
    }

    /**
     * @return {string} Address of delivery
     */
    get promo() {
        return this.#promo;
    }

    /**
     * @param {string} newAddress
     */
    set address(newAddress) {
        this.#address = newAddress;
    }

    /**
     * @param {Object} newRecipient
     */
    set recipient(newRecipient) {
        this.#recipient = newRecipient;
    }

    /**
     * @param {String} newPromo
     */
    set promo(newPromo) {
        this.#promo = newPromo;
    }

    /**
     * @description Loads all information about order via AJAX
     */
    loadOrder = () => {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.order,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw Responses.Error;
            }
            return response.json();
        }).then((parsedJson) => {
            this.#price = parsedJson.price;
            this.#products = parsedJson.products;
            this.#recipient = parsedJson.recipient;
            ymaps.geolocation.get()
                .then((parsedJson) => {
                    ymaps.geocode(parsedJson.geoObjects.position, {'json': true})
                        .then((parsedJson) => {
                            this.address = parsedJson.GeoObjectCollection
                                .featureMember[0]
                                .GeoObject
                                .metaDataProperty
                                .GeocoderMetaData
                                .Address
                                .formatted;
                            this.bus.emit(Events.OrderLoaded, Responses.Success);
                        });
                }).catch(() => {
                    this.bus.emit(Events.OrderLoaded, Responses.Success);
                });
        }).catch(() => {
            this.bus.emit(Events.OrderLoaded, Responses.Error);
        });
    }

    /**
     * @description Sends order form to backend
     */
    sendOrder = () => {
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.order,
            body: {address: {address: this.#address},
                recipient: this.#recipient,
                promo_code: this.promo ? this.promo : '',
            },
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Offline: {
                Router.open('/offline');
                break;
            }
            default: {
                console.error('error order sending');
            }
            }
        });
        this.promo = null;
    }

    /**
     * @description Sends promo
     */
    sendPromo = () => {
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.promo,
            body: {
                products: this.products.map((product) => {
                    return product.id;
                }),
                promo_code: this.promo,
            },
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((response) => {
            this.bus.emit(Events.PromoSent, Responses.Success, response);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.PromocodeNoPromoProducts: {
                this.bus.emit(Events.PromoSent, Responses.PromocodeNoPromoProducts);
                break;
            }
            case HTTPResponses.InternalError:
            case HTTPResponses.PromocodeIncorrectPromo: {
                this.bus.emit(Events.PromoSent, Responses.PromocodeIncorrectPromo);
                break;
            }
            case HTTPResponses.Offline: {
                Router.open('/offline');
                break;
            }
            default: {
                this.bus.emit(Events.PromoSent, Responses.Error);
            }
            }
        });
    }
}

export default OrderModel;
