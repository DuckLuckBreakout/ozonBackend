import BasePresenter from './BasePresenter.js';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import Router from '../utils/router/Router';
import Bus from '../utils/bus/bus';

/**
 * @description Presenter for Login View and Model
 */
class OrderPresenter extends BasePresenter {
    /**
     * @param {HTMLElement} application html of application
     * @param {Class} View Class of view object
     * @param {Class} Model Class of model object
     */
    constructor(application, View, Model) {
        super(application, View, Model);
        this.bus.on(Events.OrderLoad, this.loadOrder);
        this.bus.on(Events.OrderLoaded, this.orderLoadedReaction);
        this.bus.on(Events.SendOrder, this.sendOrder);
        this.bus.on(Events.SendPromo, this.sendPromo);
        this.bus.on(Events.PromoSent, this.promoSentReaction);
    }

    /**
     * @return {Array} array of products
     */
    get products() {
        return this.model.products;
    }

    /**
     * @return {Object} total price, discount price, discount size
     */
    get price() {
        return this.model.price;
    }

    /**
     * @return {Object} Information about recipient
     */
    get recipient() {
        return this.model.recipient;
    }

    /**
     * @return {String} Address
     */
    get address() {
        return this.model.address;
    }

    /**
     * @return {String} Promo
     */
    get promo() {
        return this.model.promo;
    }

    /**
     * @param {String} newPromo
     */
    set promo(newPromo) {
        this.model.promo = newPromo;
    }

    /**
     * @description Loads all information about order from model
     */

    loadOrder = () => {
        this.model.loadOrder();
    }

    /**
     * @param {string} result
     * @description Reaction on order loading
     */
    orderLoadedReaction = (result) => {
        if (result === Responses.Success) {
            this.view.render();
            this.view.cache.hidden = false;
        } else {
            console.error('Cant load order');
            Router.open('/items', {replaceState: true});
        }
    }

    /**
     * @param {string} result
     * @param {Object|null} newBillInfo
     * @description Reaction on promo sent
     */
    promoSentReaction = (result, newBillInfo = null) => {
        switch (result) {
        case Responses.Success: {
            this.view.drawNewBill(newBillInfo, this.promo);
            break;
        }
        case Responses.PromocodeNoPromoProducts: {
            this.view.drawNoPromoProducts(this.promo);
            break;
        }
        case Responses.PromocodeIncorrectPromo: {
            this.view.drawIncorrectPromo(this.promo);
            break;
        }
        }
    }

    /**
     * @description Sends order form to model
     */
    sendOrder = () => {
        this.model.address = document.getElementsByName('address')[0].value.trim();
        const email = document.getElementsByName('email')[0].value.trim();
        const first_name = document.getElementsByName('firstName')[0].value.trim();
        const last_name = document.getElementsByName('lastName')[0].value.trim();

        this.model.recipient = {first_name, last_name, email};
        this.model.sendOrder();
        Bus.globalBus.emit(Events.OrderSent);
    }

    /**
     * @description Sends order form to model
     */
    sendPromo = () => {
        this.promo = document.getElementsByName('promo')[0].value;
        this.model.sendPromo();
    }
}

export default OrderPresenter;
