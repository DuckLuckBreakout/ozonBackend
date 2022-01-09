import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import Bus from '../utils/bus/bus.js';
import HTTPResponses from '../utils/http-responses/httpResponses';

/**
 * @description Model for Header in MVP Arch
 */
class HeaderModel extends BaseModel {
    #categories;
    #currentCategoryIndex = 0;
    #userIsAuthorized = false;

    /**
     * @return {Object} categories
     */
    get categories() {
        return this.#categories;
    }

    /**
     * @return {Number} current category index
     */
    get currentCategoryIndex() {
        return this.#currentCategoryIndex;
    }

    /**
     * swap user authorize state
     */
    swapAuthorizeState() {
        this.#userIsAuthorized = !this.#userIsAuthorized;
    }

    /**
     * @return {boolean} current authorize state
     */
    get currentAuthorizeState() {
        return this.#userIsAuthorized;
    }

    /**
     * @param {Number} index
     */
    set currentCategoryIndex(index) {
        this.#currentCategoryIndex = index;
    }

    /**
     *
     * @description Loads item IN MODEL!
     */
    loadHeader() {
        AjaxModule.getUsingFetch({
            url: serverApiPath + '/category',
        }).then((response) => {
            return response.json();
        }).then((parsedJson) => {
            this.#categories = parsedJson;
            Bus.globalBus.emit(Events.HeaderChangeCatalogID, parsedJson[0].id);
            this.bus.emit(Events.HeaderLoaded, Responses.Success);
        }).catch(() => {
            this.bus.emit(Events.HeaderLoaded, Responses.Error);
        });
    }

    /**
     * @description check user authorization
     */
    checkAuthorizeState() {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.sessionUrl,
            body: null,
        }).then((response) => {
            if (response.status === HTTPResponses.Success) {
                Bus.globalBus.emit(Events.LoginEmitResult, Responses.Success);
            } else {
                Bus.globalBus.emit(Events.LoginEmitResult, Responses.Error);
            }
        }).catch(() => {
            Bus.globalBus.emit(Events.LoginEmitResult, Responses.Error);
        });
    }
}

export default HeaderModel;
