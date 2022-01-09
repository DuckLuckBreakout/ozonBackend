import BaseView from '../BaseView.js';
import headerStyles from './HeaderView.scss';
import buttonStyles from '../Common/Button/Button.scss';
import imgStyles from '../Common/Img/Img.scss';
import inputStyles from '../Common/Input/Input.scss';
import textStyles from '../Common/TextArea/TextArea.scss';
import linkStyles from '../Common/Link/Link.scss';
import decorators from '../decorators.scss';
import headerTemplate from './HeaderView.hbs';
import Events from '../../utils/bus/events';
import Img from '../Common/Img/Img';
import Button from '../Common/Button/Button';
import Input from '../Common/Input/Input';
import Router from '../../utils/router/Router';
import Bus from '../../utils/bus/bus';

/**
 * @class HeaderView
 * @extends BaseView
 * @classdesc Class for Header page
 */
class HeaderView extends BaseView {
    /**
     * @param {Object} URLParams
     */
    show = (URLParams = {}) => {
        this.bus.emit(Events.HeaderLoad, this.ID);
    }

    /**
     *
     */
    render = () => {
        this.parent.innerHTML = '';
        const logo = new Img({src: '/svg/header/logo.svg'});
        const headerMenu = [
            {
                text: 'Войти',
                img: new Img({src: '/svg/header/smile.svg'}),
                button: new Button(),
                href: '/login',
                name: 'login',
            }, {
                text: 'Заказы',
                img: new Img({src: '/svg/header/orders.svg'}),
                button: new Button(),
                href: '/orders',
                name: 'orders',
            }, {
                text: 'Корзина',
                img: new Img({src: '/svg/header/cart.svg'}),
                button: new Button(),
                href: '/cart',
                name: 'cart',
                additionalSpan: true,
                additionalSpanClass: headerStyles.menuItemCounter,
            },
        ];
        const catalog = {
            text: 'Каталог',
            img: [
                new Img({src: '/svg/header/catalog.svg'}),
                new Img({src: '/svg/header/catalog_close.svg'}),
            ],
            categories: this.presenter.categories,
        };
        const searchBlock = {
            searchButton: {
                button: new Button(),
                img: new Img({src: '/svg/header/search.svg'}),
            },
            searchInput: new Input({
                placeholder: 'Искать на Ozon',
                type: 'text',
                name: 'search',
            }),
            searchArea: {
                text: 'Везде',
                img: new Img({src: '/svg/header/area.svg'}),
            },
        };
        const template = headerTemplate({
            logo: logo,
            headerMenu: headerMenu,
            catalog: catalog,
            searchBlock: searchBlock,
            currentCategory: this.ID,
            headerStyles: headerStyles,
            imgStyles: imgStyles,
            inputStyles: inputStyles,
            buttonStyles: buttonStyles,
            linkStyles: linkStyles,
            textStyles: textStyles,
            decorators: decorators,
        });

        this.cache = new DOMParser().parseFromString(template, 'text/html')
            .getElementsByClassName(headerStyles.main)[0];

        const logoButton = this.cache.getElementsByClassName(headerStyles.logoBlock)[0];
        logoButton.addEventListener('click', (evt) => {
            evt.preventDefault();
            this.dropSearchInput();
            Router.open('/', {dropFilter: true, dropSort: true});
        });

        const catalogBlock = this.cache.getElementsByClassName(headerStyles.catalogBlock)[0];
        const catalogList = this.cache.getElementsByClassName(headerStyles.catalogList)[0];
        const catalogLeftBackground = this.cache.getElementsByClassName(headerStyles.catalogLeftBackground)[0];
        const catalogRightBackground = this.cache.getElementsByClassName(headerStyles.catalogRightBackground)[0];
        catalogBlock.addEventListener('click', () => {
            const images = Array.from(catalogBlock.getElementsByClassName(imgStyles.headerMenu));
            if (images[0].classList.contains(decorators.hidden)) {
                images[0].classList.remove(decorators.hidden);
                images[1].classList.add(decorators.hidden);
                catalogList.classList.add(decorators.hidden);
                catalogLeftBackground.classList.add(decorators.hidden);
                catalogRightBackground.classList.add(decorators.hidden);
            } else {
                images[1].classList.remove(decorators.hidden);
                images[0].classList.add(decorators.hidden);
                catalogList.classList.remove(decorators.hidden);
                catalogLeftBackground.classList.remove(decorators.hidden);
                const catalogListHeight = catalogList.getBoundingClientRect().height + 40;
                catalogLeftBackground.style.height = `${catalogListHeight}px`;
                catalogRightBackground.classList.remove(decorators.hidden);
                catalogRightBackground.style.height = `${catalogListHeight}px`;
                catalogList.style.height = `${catalogListHeight}px`;
            }
        });

        this.cache.addEventListener('click', (evt) => {
            if (evt.target.hasAttribute('category')) {
                evt.preventDefault();
                const categoryId = parseInt(evt.target.getAttribute('category'));
                catalogBlock.dispatchEvent(new Event('click'));
                this.dropSearchInput();
                Router.open(`/items/${categoryId}`, {dropFilter: true, dropSort: true});
            }
        });

        const menuItems = Array.from(this.cache.getElementsByClassName(headerStyles.menuItem));
        menuItems.forEach((menuItem) => {
            menuItem.addEventListener('click', (event) => {
                event.preventDefault();
                const href = menuItem.getAttribute('href');
                Router.open(href);
            });
        });

        const searchForm = this.cache.getElementsByClassName(headerStyles.searchForm)[0];
        searchForm.addEventListener('submit', (event) => {
            event.preventDefault();
            const searchInput = this.cache.getElementsByClassName(inputStyles.search)[0];
            Router.open('/search/1/', {dropFilter: true, dropSort: true}, {text: searchInput.value});
        });

        const catalogListCategories = Array.from(this.cache.getElementsByClassName(headerStyles.category));
        const catalogListSubcategoriesWrapper = Array.from(
            this.cache.getElementsByClassName(headerStyles.catalogListSubcategoriesWrapper),
        );
        catalogListCategories.forEach((category, i) => {
            category.addEventListener('mouseover', () => {
                if (i !== this.presenter.currentCategoryIndex) {
                    category.style.backgroundColor = 'white';
                    category.getElementsByClassName(linkStyles.catalogCategory)[0].style.color = '#005bff';
                    catalogListSubcategoriesWrapper[i].classList.remove(decorators.hidden);
                    catalogListSubcategoriesWrapper[this.presenter.currentCategoryIndex]
                        .classList.add(decorators.hidden);
                    catalogListCategories[this.presenter.currentCategoryIndex].removeAttribute('style');
                    catalogListCategories[this.presenter.currentCategoryIndex]
                        .getElementsByClassName(linkStyles.catalogCategory)[0].style.color = '#001a34';
                    this.presenter.currentCategoryIndex = i;
                }
            });
        });
        this.parent.appendChild(this.cache);
        Bus.globalBus.emit(Events.CartLoadProductsAmount);
    }

    /**
     *
     * @param {boolean} userIsAuthorized
     */
    changeLoginIcon = (userIsAuthorized) => {
        const icon = document.getElementsByName('login')[0];
        if (userIsAuthorized === true) {
            icon.getElementsByTagName('span')[0].textContent = 'Профиль';
            icon.setAttribute('href', '/profile');
        } else {
            icon.getElementsByTagName('span')[0].textContent = 'Войти';
            icon.setAttribute('href', '/login');
        }
    }

    dropSearchInput = () => {
        document.getElementsByClassName(inputStyles.search)[0].value = '';
    }

    /**
     * @param {number} value
     */
    changeCartItems = (value) => {
        const counter = document.getElementsByClassName(headerStyles.menuItemCounter)[0];
        const newCounterAmount = +counter.innerHTML + value;
        counter.innerHTML = newCounterAmount.toString();
        counter.hidden = !newCounterAmount;
    }

    /**
     * @param {number} value
     * @description sets items amount in cart to 0
     */
    setCartItems = (value) => {
        const counter = document.getElementsByClassName(headerStyles.menuItemCounter)[0];
        counter.innerHTML = value.toString();
        counter.hidden = !value;
    }
}

export default HeaderView;
