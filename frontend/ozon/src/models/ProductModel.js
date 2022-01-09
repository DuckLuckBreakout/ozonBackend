import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';

/**
 * @description Model for Product in MVP Arch
 */
class ProductModel extends BaseModel {
    #item
    #recommendations
    #sortKey
    #sortDirection

    /**
     * @return {Object} item
     */
    get item() {
        return this.#item;
    }

    /**
     * @return {string}
     */
    get sortKey() {
        return this.#sortKey;
    }

    /**
     * @return {string}
     */
    get sortDirection() {
        return this.#sortDirection;
    }

    /**
     * @return {Array}
     */
    get recommendations() {
        return this.#recommendations;
    }

    /**
     * @param {string} newSortKey
     */
    set sortKey(newSortKey) {
        this.#sortKey = newSortKey;
    }

    /**
     * @param {string} newSortDirection
     */
    set sortDirection(newSortDirection) {
        this.#sortDirection = newSortDirection;
    }

    /**
     * @param {Array} newRecommendations
     */
    set recommendations(newRecommendations) {
        this.#recommendations = newRecommendations;
    }

    /**
     *
     * @param {Number} itemId ID of item on server
     * @description Loads item IN MODEL!
     */
    loadProduct(itemId) {
        AjaxModule.getUsingFetch({
            url: serverApiPath + `/product/${itemId}`,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            const base = parsedJson['price']['base_cost'];
            const discount = parsedJson['price']['discount'];
            const discountPrice = Math.ceil(base * (1 - discount * 0.01));
            this.#item = {
                name: parsedJson['title'],
                price: {
                    discountPrice: discountPrice,
                    base: base,
                    discount: discount,
                },
                rating: parsedJson['rating'],
                description: {
                    descriptionText: parsedJson['description'],
                    category: parsedJson['category'],
                },
                properties: JSON.parse(parsedJson['properties']),
                images: parsedJson['images'],
                id: parsedJson['id'],
                category_path: parsedJson['category_path'],
            };
            this.bus.emit(Events.ProductLoaded, Responses.Success);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.ProductLoaded, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.ProductLoaded, Responses.Offline);
                break;
            }
            default: {
                this.bus.emit(Events.ProductLoaded, Responses.Error);
                break;
            }
            }
        });
    }

    loadStarsCounter = () => {
        AjaxModule.getUsingFetch({
            url: serverApiPath + `/review/statistics/product/${this.item.id}`,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((response) => {
            Object.assign(this.item, {stars: response.stars});
            this.bus.emit(Events.ProductStarsCounterLoaded, Responses.Success);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.ProductStarsCounterLoaded, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.ProductStarsCounterLoaded, Responses.Offline);
                break;
            }
            default: {
                this.bus.emit(Events.ProductStarsCounterLoaded, Responses.Error);
                break;
            }
            }
        });
    }

    /**
     *
     * @param {Number} itemId ID of item on server
     * @description Loads recommendation IN MODEL!
     */
    loadRecommendations(itemId) {
        AjaxModule.postUsingFetch({
            url: serverApiPath + `/product/recommendations/${itemId}`,
            body: {
                count: 10,
            },
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            this.recommendations = parsedJson;
            this.bus.emit(Events.RecommendationLoaded, Responses.Success);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.RecommendationLoaded, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.bus.emit(Events.RecommendationLoaded, Responses.Offline);
                break;
            }
            default: {
                this.bus.emit(Events.RecommendationLoaded, Responses.Error);
                break;
            }
            }
        });
    }
}

export default ProductModel;
