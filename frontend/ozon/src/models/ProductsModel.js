import AjaxModule from '../modules/Ajax/Ajax';
import {serverApiPath} from '../utils/urls/urls';
import BaseModel from './BaseModel';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import HTTPResponses from '../utils/http-responses/httpResponses';

/**
 * @description Model for ProductS in MVP Arch
 */
class ProductsModel extends BaseModel {
    #products
    #paginationInfo
    #categoryId
    #sortKey
    #sortDirection
    #filter

    /**
     *
     * @return {Object} Array of products
     */
    get products() {
        return this.#products;
    }

    /**
     *
     * @return {Object} Pagination info
     */
    get paginationInfo() {
        return this.#paginationInfo;
    }

    /**
     * @return {Number}
     */
    get categoryId() {
        return this.#categoryId;
    }

    /**
     * @return {String}
     */
    get sortKey() {
        return this.#sortKey;
    }

    /**
     * @return {String}
     */
    get sortDirection() {
        return this.#sortDirection;
    }

    /**
     * @return {Object}
     */
    get filter() {
        return this.#filter;
    }

    /**
     *
     * @param {String} sortKey
     */
    set sortKey(sortKey) {
        this.#sortKey = sortKey;
    }

    /**
     *
     * @param {String} sortDirection
     */
    set sortDirection(sortDirection) {
        this.#sortDirection = sortDirection;
    }

    /**
     *
     * @param {Number} id
     */
    set categoryId(id) {
        this.#categoryId = id;
    }

    /**
     *
     * @param {Object} newProducts
     */
    set products(newProducts) {
        return this.#products = newProducts;
    }

    /**
     * @param {Object} newFilter
     */
    set filter(newFilter) {
        this.#filter = newFilter;
    }

    /**
     * @param {Number|String} category
     * @param {Number|String} page
     * @param {String} sortKey
     * @param {String} sortDirection
     * @param {Object} body Body of request
     */
    loadProducts(category, page, sortKey, sortDirection, body = {
        page_num: +page,
        count: 9,
        sort_key: sortKey,
        sort_direction: sortDirection,
        category: +category,
    }) {
        if (this.filter !== undefined) {
            body['filter'] = {...this.filter};
            body.filter.min_price = body.filter.min_price === undefined ? 0: body.filter.min_price;
            body.filter.max_price = body.filter.max_price === undefined ? 1e6: body.filter.max_price;
        }
        AjaxModule.postUsingFetch({
            url: serverApiPath + '/product',
            body: body,
        }).then((response) => {
            if (!response.ok) {
                throw response.status;
            }

            return response.json();
        }).then((parsedJson) => {
            this.#products = parsedJson['list_preview_products'];
            this.#paginationInfo = {
                pagesCount: parsedJson['max_count_pages'],
                currentPage: page,
            };
            this.bus.emit(Events.ProductsLoaded, Responses.Success);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.ProductsLoaded, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.#products = [];
                this.#paginationInfo = {
                    pagesCount: 1,
                    currentPage: 1,
                };
                this.bus.emit(Events.ProductsLoaded, Responses.Success);
                break;
            }
            default: {
                this.bus.emit(Events.ProductsLoaded, Responses.Error);
                break;
            }
            }
        });
    }


    /**
     * @param {String} searchData
     * @param {Number|String} page
     * @param {String} sortKey
     * @param {String} sortDirection
     * @param {Object} body Body of request
     */
    loadProductsSearch(searchData, page, sortKey, sortDirection, body = {
        query_string: searchData,
        page_num: +page,
        count: 9,
        sort_key: sortKey,
        sort_direction: sortDirection,
        category: 1,
    }) {
        if (this.filter !== undefined) {
            body['filter'] = {...this.filter};
            body.filter.min_price = body.filter.min_price === undefined ? 0: body.filter.min_price;
            body.filter.max_price = body.filter.max_price === undefined ? 1e6: body.filter.max_price;
        }
        AjaxModule.postUsingFetch({
            url: serverApiPath + '/product/search',
            body: body,
        }).then((response) => {
            if (response.status !== HTTPResponses.Success) {
                throw response.status;
            }
            return response.json();
        }).then((parsedJson) => {
            this.#products = parsedJson['list_preview_products'];
            this.#paginationInfo = {
                pagesCount: parsedJson['max_count_pages'],
                currentPage: page,
            };
            this.bus.emit(Events.ProductsLoaded, Responses.Success);
        }).catch((err) => {
            switch (err) {
            case HTTPResponses.Unauthorized: {
                this.bus.emit(Events.ProductsLoaded, Responses.Unauthorized);
                break;
            }
            case HTTPResponses.Offline: {
                this.#products = [];
                this.#paginationInfo = {
                    pagesCount: 1,
                    currentPage: 1,
                };
                this.bus.emit(Events.ProductsLoaded, Responses.Success);
                break;
            }
            default: {
                this.bus.emit(Events.ProductsLoaded, Responses.Error);
                break;
            }
            }
        });
    }
}

export default ProductsModel;
