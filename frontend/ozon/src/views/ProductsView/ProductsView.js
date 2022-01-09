import BaseView from '../BaseView.js';
import ListOfProducts from '../Common/ListOfProducts/ListOfProducts.js';
import ListOfProductsItemStyles from '../Common/ListOfProducts/ListOfProductsItem/ListOfProductsItem.scss';
import Pagination from '../Common/Pagination/Pagination';
import productsPageTemplate from './ProductsView.hbs';
import productsStyles from './ProductsView.scss';
import buttonStyles from './../Common/Button/Button.scss';
import textStyles from './../Common/TextArea/TextArea.scss';
import Filter from '../Common/Filter/Filter.js';
import Bus from '../../utils/bus/bus';
import Router from '../../utils/router/Router';
import Events from '../../utils/bus/events';

/**
 * @class ProductsView
 * @extends BaseView
 * @classdesc Class for showing product
 */
class ProductsView extends BaseView {
    static #types = {
        search: 'search',
        category: 'category',
    }

    #viewType
    #filterHTML

    /**
     * @param {Object} URLParams
     */
    show = (URLParams = {}) => {
        if (!this.IDs) {
            this.IDs = {};
        }
        if (!this.IDs['page']) {
            this.IDs['page'] = 1;
        }
        if (!this.IDs['category']) {
            this.IDs['category'] = 1;
        }

        if (!this.presenter.sortKey) {
            this.presenter.changeSortKey('cost');
        }
        if (!this.presenter.sortDirection) {
            this.presenter.changeSortDirection('ASC');
        }

        if (Object.keys(URLParams).length > 1) {
            this.presenter.changeSortKey(URLParams.sortKey);
            this.presenter.changeSortDirection(URLParams.sortDirection);
            this.presenter.setFilter(URLParams);
        } else {
            this.presenter.dropFilter();
            this.dropFilter();
            this.presenter.dropSort();
            this.dropSort();
        }

        if (this.IDs['dropFilter']) {
            this.presenter.dropFilter();
            this.dropFilter();
        }

        if (this.IDs['dropSort']) {
            this.presenter.dropSort();
            this.dropSort();
        }

        if (URLParams && URLParams.text) {
            this.#viewType = ProductsView.#types.search;
            this.IDs['searchText'] = URLParams.text;
            this.bus.emit(
                Events.ProductsLoadSearch,
                URLParams.text,
                this.IDs['page'],
                this.presenter.sortKey,
                this.presenter.sortDirection,
            );
        } else {
            this.#viewType = ProductsView.#types.category;
            this.bus.emit(
                Events.ProductsLoad,
                this.IDs['category'],
                this.IDs['page'],
                this.presenter.sortKey,
                this.presenter.sortDirection,
            );
        }
    }

    render = () => {
        this.goUp();

        this.parent.innerHTML = '';
        const productsListHtmlString = new ListOfProducts(this.presenter.products, 'products').getHtmlString();
        const pagination = new Pagination(this.presenter.paginationInfo).getHtmlString();
        const filter = this.#filterHTML ? this.#filterHTML : Filter.getHtmlString();
        const template = productsPageTemplate({
            isEmpty: this.presenter.products.length === 0,
            productsList: productsListHtmlString,
            pagination: pagination,
            filter: filter,
            select: [
                {
                    key: 'cost',
                    direction: 'ASC',
                    name: 'Сначала дешевые',
                },
                {
                    key: 'cost',
                    direction: 'DESC',
                    name: 'Сначала дорогие',
                },
                {
                    key: 'date',
                    direction: 'DESC',
                    name: 'Новинки',
                },
                {
                    key: 'rating',
                    direction: 'DESC',
                    name: 'Высокий рейтинг',
                },
                {
                    key: 'discount',
                    direction: 'DESC',
                    name: 'По размеру скидки',
                },
            ],
            sort: {
                key: this.presenter.sortKey,
                direction: this.presenter.sortDirection,
            },
            productsStyles: productsStyles,
            buttonStyles: buttonStyles,
            textStyles: textStyles,
        });
        this.cache = new DOMParser().parseFromString(template, 'text/html')
            .getElementsByClassName(productsStyles.block)[0];

        const select = this.cache.getElementsByClassName(productsStyles.select)[0];
        select.addEventListener('change', () => {
            const selected = select.selectedOptions[0];
            const sortKey = selected.getAttribute('key');
            const sortDirection = selected.getAttribute('direction');
            this.presenter.changeSortKey(sortKey);
            this.presenter.changeSortDirection(sortDirection);
            this.#viewType === ProductsView.#types.category ?
                Router.open(`/items/${this.IDs['category']}`, {}, this.presenter.getParams()) :
                Router.open('/search/1/',
                    {},
                    Object.assign(this.presenter.getParams(),
                        {text: this.IDs['searchText']}));
        });

        for (const button of this.cache.getElementsByClassName(buttonStyles.paginator)) {
            button.addEventListener('click', () => {
                const page = parseInt(button.getAttribute('page'));
                this.ID = page;
                this.#viewType === ProductsView.#types.category ?
                    Router.open(`/items/${this.IDs['category']}/${page}`, {}, this.presenter.getParams()) :
                    Router.open(`/search/${page}/`,
                        {},
                        Object.assign(this.presenter.getParams(),
                            {text: this.IDs['searchText']}));
            });
        }


        for (const itemContainer of this.cache.getElementsByClassName(ListOfProductsItemStyles.block)) {
            const productID = parseInt(itemContainer.getAttribute('item-id'));
            itemContainer.getElementsByClassName(ListOfProductsItemStyles.infoWrapper)[0]
                .addEventListener('click', () => {
                    Bus.globalBus.emit(Events.ProductChangeID, productID);
                    Router.open(`/item/${productID}`);
                });

            let item = itemContainer.getElementsByClassName(buttonStyles.notInCartProducts)[0];
            if (item === undefined) {
                item = itemContainer.getElementsByClassName(buttonStyles.inCartProducts)[0];
            }

            item.addEventListener('click', (evt) => {
                evt.preventDefault();
                Bus.globalBus.emit(Events.CartAddProduct, productID, 1, Events.ProductsItemAdded);
            });
        }

        this.parent.appendChild(this.cache);
        Bus.globalBus.emit(Events.CartGetProductsID, Events.ProductsCartLoadedProductsID);

        document.getElementById('filtration_form').addEventListener('submit', (event) => {
            event.preventDefault();
            this.IDs['page'] = 1;
            this.#viewType === ProductsView.#types.category ?
                Router.open(`/items/${this.IDs['category']}/1`, {}, this.presenter.getParams()) :
                Router.open('/search/1/',
                    {},
                    Object.assign(this.presenter.getParams(),
                        {text: this.IDs['searchText']}));
        });

        document.getElementById('drop-filters').addEventListener('click', (event) => {
            event.preventDefault();
            this.presenter.dropFilter();
            this.dropFilter();
            this.#viewType === ProductsView.#types.search ? this.show({text: this.IDs['searchText']}) : this.show();
        });

        const min_price_input = document.getElementById('min_price');
        const max_price_input = document.getElementById('max_price');
        min_price_input.addEventListener('keypress', (event) => {
            if (event.key === '-') {
                event.preventDefault();
            }
        });

        max_price_input.addEventListener('keypress', (event) => {
            if (event.key === '-') {
                event.preventDefault();
            }
        });

        const onMainElem = this.cache.getElementsByClassName(productsStyles.onMainPage)[0];
        if (onMainElem) {
            onMainElem.addEventListener('click', (event) => {
                event.preventDefault();
                Router.open('/', {dropFilter: true, dropSort: true});
            });
        }
        this.drawFilter();
    };


    /**
     * @param {Set} productsInCart
     */
    setAddedProducts = (productsInCart) => {
        for (const item of document.getElementsByClassName(ListOfProductsItemStyles.block)) {
            if (productsInCart.has(+item.getAttribute('item-id'))) {
                this.setButtonAdded(item.getElementsByTagName('button')[0], +item.getAttribute('item-id'));
            }
        }
    }

    /**
     * @param {number} productID
     */
    setProductAdded = (productID) => {
        const item = Array.from(document.getElementsByClassName(ListOfProductsItemStyles.block)).filter((item) => {
            return +item.getAttribute('item-id') === productID;
        })[0];
        if (item) {
            this.setButtonAdded(item.getElementsByClassName(
                ListOfProductsItemStyles.buttonsBlock)[0].getElementsByTagName('button')[0],
            productID);
        }
    }

    /**
     * @param {HTMLElement} button
     * @param {number} id
     */
    setButtonAdded = (button, id) => {
        const newButton = button.cloneNode(true);
        this.setButtonAddedStyle(newButton);
        newButton.addEventListener('click', () => {
            Bus.globalBus.emit(Events.CartRemoveProduct, id, Events.ProductsItemNotAdded);
        });
        button.replaceWith(newButton);
    }

    /**
     * @param {number} productID
     */
    setProductNotAdded = (productID) => {
        const item = Array.from(document.getElementsByClassName(ListOfProductsItemStyles.block)).filter((item) => {
            return +item.getAttribute('item-id') === +productID;
        })[0];
        if (item) {
            const button = item.getElementsByTagName('button')[0];
            const newButton = button.cloneNode(true);
            this.setButtonNotAddedStyle(newButton);
            newButton.addEventListener('click', () => {
                const id = item.getAttribute('item-id');
                Bus.globalBus.emit(Events.CartAddProduct, +item.getAttribute('item-id'), 1, Events.ProductsItemAdded);
                this.setProductAdded(+id);
            });
            button.replaceWith(newButton);
        }
    }

    /**
     * @param {HTMLElement} button
     */
    setButtonAddedStyle = (button) => {
        button.className = buttonStyles.inCartProducts;
        button.getElementsByTagName('span')[0].innerHTML = 'В корзине';
    }

    /**
     * @param {HTMLElement} button
     */
    setButtonNotAddedStyle = (button) => {
        button.className = buttonStyles.notInCartProducts;
        button.getElementsByTagName('span')[0].innerHTML = 'В корзину';
    }

    drawFilter = () => {
        const filter = this.presenter.filter;
        if (!filter) {
            return;
        }
        document.getElementById('min_price').value = filter.min_price !== undefined ? filter.min_price : '';
        document.getElementById('max_price').value = filter.max_price !== undefined ? filter.max_price : '';
        document.getElementById('is_new').checked = filter.is_new;
        document.getElementById('is_rating').checked = filter.is_rating;
        document.getElementById('is_discount').checked = filter.is_discount;
    }

    dropFilter = () => {
        const form = document.getElementById('filtration_form');
        if (!form) {
            return;
        }

        for (const input of form.getElementsByTagName('input')) {
            input.type === 'checkbox' ? input.checked = false : input.value = '';
        }
    }

    dropSort = () => {
        const sort = document.getElementsByClassName(productsStyles.select)[0];
        if (!sort) {
            return;
        }

        sort.key = 'cost';
        sort.direction = 'ASC';
        sort.value = 'Сначала дешевые';
    }

    drawIncorrectFilterWarning = () => {
        const incorrectFilterLabel = document.getElementById('incorrect_filter_label');
        if (!incorrectFilterLabel) {
            return;
        }
        incorrectFilterLabel.innerHTML = 'Некорректные данные для фильтрации!';
    }

    dropIncorrectFilterWarning = () => {
        const incorrectFilterLabel = document.getElementById('incorrect_filter_label');
        if (!incorrectFilterLabel) {
            return;
        }
        incorrectFilterLabel.innerHTML = '';
    }
}

export default ProductsView;
