import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath, urls} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';
import Router from '../utils/router/Router';
import Bus from '../utils/bus/bus';

/**
 * @description Model for Review in MVP Arch
 */
class ReviewModel extends BaseModel {
    #product
    #rating
    #advantages
    #disadvantages
    #comment
    #isPublic
    #userName
    #reviews
    #paginationInfo

    /**
     * @return {Object} product
     */
    get product() {
        return this.#product;
    }

    /**
     * @return {Number} rating
     */
    get rating() {
        return this.#rating;
    }

    /**
     * @return {string} advantages
     */
    get advantages() {
        return this.#advantages;
    }

    /**
     * @return {string} disadvantages
     */
    get disadvantages() {
        return this.#disadvantages;
    }

    /**
     * @return {string} comment
     */
    get comment() {
        return this.#comment;
    }

    /**
     * @return {boolean} isPublic
     */
    get isPublic() {
        return this.#isPublic;
    }

    /**
     * @return {String} userName
     */
    get userName() {
        return this.#userName;
    }

    /**
     * @return {Array}
     */
    get reviews() {
        return this.#reviews;
    }

    /**
     * @return {Object}
     */
    get paginationInfo() {
        return this.#paginationInfo;
    }


    /**
     * @param {Object} product
     */
    set product(product) {
        this.#product = product;
    }

    /**
     * @param {Number} rating
     */
    set rating(rating) {
        this.#rating = rating;
    }

    /**
     * @param {string} advantages
     */
    set advantages(advantages) {
        this.#advantages = advantages;
    }

    /**
     * @param {string} disadvantages
     */
    set disadvantages(disadvantages) {
        this.#disadvantages = disadvantages;
    }

    /**
     * @param {string} comment
     */
    set comment(comment) {
        this.#comment = comment;
    }

    /**
     * @param {boolean} isPublic
     */
    set isPublic(isPublic) {
        this.#isPublic = isPublic;
    }

    /**
     * @param {String} userName
     */
    set userName(userName) {
        this.#userName = userName;
    }

    /**
     * @param {Array} newReviews
     */
    set reviews(newReviews) {
        this.#reviews = newReviews;
    }

    /**
     * @param {Object} newPaginationInfo
     */
    set paginationInfo(newPaginationInfo) {
        this.#paginationInfo = newPaginationInfo;
    }

    /**
     * @description Loads all information about review via AJAX
     */
    loadReviewRights = () => {
        AjaxModule.getUsingFetch({
            url: serverApiPath + urls.reviewRights + this.product.id,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            this.bus.emit(Events.ReviewRightsLoaded, Responses.Success);
        }).catch(() => {
            this.bus.emit(Events.ReviewRightsLoaded, Responses.Error);
        });
    }

    /**
     * @description Sends review form to backend
     */
    sendReview = () => {
        const body = {
            product_id: this.#product.id,
            rating: this.#rating,
            advantages: this.#advantages,
            disadvantages: this.#disadvantages,
            comment: this.#comment,
            is_public: this.#isPublic !== undefined ? this.#isPublic : true,
        };
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.review,
            body: body,
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
                console.error('error review sending');
            }
            }
        });
    }

    /**
     * @param {number} productID
     * @param {number} page
     * @param {string} sortKey
     * @param {string} sortDirection
     * @param {Object} body
     */
    getProductReviews = (productID, page = 1, sortKey = 'date', sortDirection = 'ASC', body = {
        page_num: +page,
        count: 1,
        sort_key: sortKey,
        sort_direction: sortDirection,
    }) => {
        AjaxModule.postUsingFetch({
            url: serverApiPath + urls.review + '/' + productID,
            body: body,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            this.reviews = parsedJson.list_reviews;
            this.paginationInfo = {
                pagesCount: parsedJson['max_count_pages'],
                currentPage: page,
            };
            Bus.globalBus.emit(Events.GetProductReviewsReaction, Responses.Success);
        }).catch(() => {
            Bus.globalBus.emit(Events.GetProductReviewsReaction, Responses.Error);
        });
    }
}

export default ReviewModel;
