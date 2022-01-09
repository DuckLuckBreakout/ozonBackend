import BasePresenter from './BasePresenter.js';
import Events from '../utils/bus/events';
import Responses from '../utils/bus/responses';
import Router from '../utils/router/Router';
import Bus from '../utils/bus/bus';

/**
 * @description Presenter for Product View and Model
 */
class ProductsPresenter extends BasePresenter {
    /**
     * @param {HTMLElement} application html of application
     * @param {Class} View Class of view object
     * @param {Class} Model Class of model object
     */
    constructor(application, View, Model) {
        super(application, View, Model);
        this.bus.on(Events.ProductsLoad, this.loadProducts);
        this.bus.on(Events.ProductsLoadSearch, this.loadSearchProducts);
        this.bus.on(Events.ProductsLoaded, this.productLoadedReaction);
        Bus.globalBus.on(Events.ProductsChangeCategory, this.changeCategory);
        Bus.globalBus.on(Events.ProductsItemAdded, this.setProductAdded);
        Bus.globalBus.on(Events.ProductsItemNotAdded, this.setProductNotAdded);
        Bus.globalBus.on(Events.HeaderChangeCategoryID, this.changeCategoryId);
        Bus.globalBus.on(Events.ProductsCartLoadedProductsID, this.productsCartGotIds);
    }

    /**
     *
     * @return {Object}
     */
    get products() {
        return this.model.products;
    }

    /**
     *
     * @return {{pagesCount: *, currentPage: Number}|*}
     */
    get paginationInfo() {
        return this.model.paginationInfo;
    }


    /**
     * @return {Number}
     */
    get categoryId() {
        return this.model.categoryId;
    }

    /**
     * @return {String}
     */
    get sortKey() {
        return this.model.sortKey;
    }

    /**
     * @return {String}
     */
    get sortDirection() {
        return this.model.sortDirection;
    }


    /**
     * @return {Object}
     */
    get filter() {
        return this.model.filter;
    }

    dropFilter = () => {
        this.model.filter = undefined;
    }

    dropSort = () => {
        this.model.sortKey = 'cost';
        this.model.sortDirection = 'ASC';
    }

    /**
     * @param {Object} URLParams
     */
    setFilter = (URLParams) => {
        /**
         * @param {string} string
         * @return {boolean}
         */
        function isStringValidBool(string) {
            if (!string) {
                return false;
            }
            string = string.toLowerCase();
            return string === 'true' || string === 'false';
        }

        this.model.filter = {};


        const priceMin = URLParams.priceMin;
        const priceMax = URLParams.priceMax;
        this.model.filter.min_price = priceMin && !isNaN(priceMin) && !isNaN(parseInt(priceMin)) && +priceMin > 0 ?
            +priceMin : undefined;
        this.model.filter.max_price = priceMax && !isNaN(priceMax) && !isNaN(parseInt(priceMax)) && +priceMax > 0?
            +priceMax : undefined;
        this.model.filter.is_new = isStringValidBool(URLParams.isNew) ?
            JSON.parse(URLParams.isNew.toLowerCase()) : false;
        this.model.filter.is_rating = isStringValidBool(URLParams.isRating) ?
            JSON.parse(URLParams.isRating.toLowerCase()) : false;
        this.model.filter.is_discount = isStringValidBool(URLParams.isDiscount) ?
            JSON.parse(URLParams.isDiscount.toLowerCase()) : false;
    }

    /**
     *
     * @param {Number} id
     */
    changeCategoryId = (id) => {
        this.model.categoryId = id;
    }

    /**
     *
     * @param {String} sortKey
     */
    changeSortKey = (sortKey) => {
        this.model.sortKey = sortKey;
    }

    /**
     *
     * @param {String} sortDirection
     */
    changeSortDirection = (sortDirection) => {
        this.model.sortDirection = sortDirection;
    }

    /**
     * @return {{priceMin: (number|undefined),
     *            sortDirection: String,
     *            priceMax: (number|undefined),
     *            sortKey: String,
     *            isRating: (boolean),
     *            isNew: (boolean),
     *            isDiscount: (boolean)}}
     */
    getParams = () => {
        this.parseFiltration();
        const filter = this.filter;
        return {
            sortKey: this.sortKey,
            sortDirection: this.sortDirection,
            priceMin: filter ? filter.min_price : undefined,
            priceMax: filter ? filter.max_price : undefined,
            isNew: filter ? filter.is_new : false,
            isRating: filter ? filter.is_rating : false,
            isDiscount: filter ? filter.is_discount: false,
        };
    }

    /**
     * @param {Number} category
     * @param {Number} page
     * @param {String} sortKey
     * @param {String} sortDirection
     */
    loadProducts = (category, page, sortKey, sortDirection) => {
        if (this.parseFiltration()) {
            this.view.dropIncorrectFilterWarning();
            this.model.loadProducts(category, page, sortKey, sortDirection);
            return;
        }
        this.view.drawIncorrectFilterWarning();
    }

    /**
     * @param {string} searchData
     * @param {number} page
     * @param {String} sortKey
     * @param {String} sortDirection
     */
    loadSearchProducts = (searchData, page, sortKey, sortDirection) => {
        if (this.parseFiltration()) {
            this.view.dropIncorrectFilterWarning();
            this.model.loadProductsSearch(searchData, page, sortKey, sortDirection);
            return;
        }
        this.view.drawIncorrectFilterWarning();
    }


    /**
     * @return {boolean}
     */
    parseFiltration = () => {
        if (!document.getElementById('min_price')) {
            return true;
        }

        const minPrice = document.getElementById('min_price').value;
        const maxPrice = document.getElementById('max_price').value;
        if (minPrice.length && maxPrice.length && +minPrice > +maxPrice) {
            return false;
        }

        this.model.filter = {
            min_price: minPrice === '' ? undefined : parseInt(minPrice),
            max_price: maxPrice === '' ? undefined : parseInt(maxPrice),
            is_new: document.getElementById('is_new').checked,
            is_rating: document.getElementById('is_rating').checked,
            is_discount: document.getElementById('is_discount').checked,
        };
        return true;
    }

    /**
     *
     * @param {string} result
     */
    productLoadedReaction = (result) => {
        switch (result) {
        case Responses.Success: {
            this.view.render();
            this.view.cache.hidden = false;
            break;
        }
        case Responses.Offline: {
            Router.open('/offline', {replaceState: true});
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
     * @param {Set} ids
     */
    productsCartGotIds = (ids) => {
        if (ids.size) {
            this.view.setAddedProducts(ids);
        }
    }

    /**
     * @param {number} newCategory
     */
    changeCategory = (newCategory) => {
        this.view.ID = newCategory;
        this.view.subID = 1; // first page of pagination!
    }

    /**
     * @param {number} itemID
     */
    setProductAdded = (itemID) => {
        this.view.setProductAdded(itemID);
    }

    /**
     * @param {number} itemID
     */
    setProductNotAdded = (itemID) => {
        this.view.setProductNotAdded(itemID);
    }
}

export default ProductsPresenter;
