import BaseView from '../BaseView.js';
import cartPageTemplate from './CartView.hbs';
import emptyCartPageTemplate from './CartViewEmpty.hbs';
import Bus from '../../utils/bus/bus';
import Events from '../../utils/bus/events';
import cartStyles from './CartView.scss';
import buttonStyles from './../Common/Button/Button.scss';
import textStyles from './../Common/TextArea/TextArea.scss';
import imgStyles from './../Common/Img/Img.scss';
import decorators from '../decorators.scss';
import Router from '../../utils/router/Router';

/**
 * @class ProductsView
 * @extends BaseView
 * @classdesc Class for showing product
 */
class CartView extends BaseView {
    /**
     * @param {Object} URLParams
     */
    show = (URLParams = {}) => {
        this.bus.emit(Events.CartLoad);
    }

    render = () => {
        this.goUp();

        this.parent.innerHTML = '';
        const price = this.presenter.price;
        let template;
        if (this.presenter.products.length) {
            template = cartPageTemplate({
                productsList: this.presenter.products,
                cartStyles: cartStyles,
                cartSize: this.presenter.products.length,
                baseCost: price.total_base_cost,
                discount: price.total_discount,
                totalCost: price.total_cost,
                decorators: decorators,
                buttonStyles: buttonStyles,
                textStyles: textStyles,
                imgStyles: imgStyles,
            });
        } else {
            template = emptyCartPageTemplate({
                cartStyles: cartStyles,
                decorators: decorators,
                textStyles: textStyles,
            });
        }
        this.cache = new DOMParser().parseFromString(template, 'text/html').getElementById('products-list-block');
        this.parent.appendChild(this.cache);


        for (const product of document.getElementsByClassName(cartStyles.productMenuText)) {
            product.addEventListener('click', (evt) => {
                evt.preventDefault();
                const counter = +document.querySelector(`div[product_id='${product.title}']`).
                    getElementsByClassName(cartStyles.count)[0].innerHTML;
                this.changeContent(+product.getAttribute('title'), -counter);
                if (!product.title) {
                    return;
                }
                Bus.globalBus.emit(Events.CartRemoveProduct, product.title, Events.ProductsItemNotAdded);
            });
        }

        for (const elemList of document.getElementsByClassName(cartStyles.productsListElem)) {
            const itemId = Number(elemList.getAttribute('product_id'));

            elemList.getElementsByClassName(buttonStyles.increment)[0].addEventListener('click', (evt) => {
                evt.preventDefault();
                const count = +elemList.getElementsByClassName(cartStyles.count)[0].textContent;
                this.changeContent(itemId, -1);
                Bus.globalBus.emit(Events.CartProductChange, {
                    id: itemId,
                    count: count - 1,
                });
            });

            elemList.getElementsByClassName(buttonStyles.decrement)[0].addEventListener('click', (evt) => {
                evt.preventDefault();
                const count = +elemList.getElementsByClassName(cartStyles.count)[0].textContent;
                this.changeContent(itemId, 1);
                Bus.globalBus.emit(Events.CartProductChange, {
                    id: itemId,
                    count: count + 1,
                });
            });
        }

        for (const itemContainer of this.cache.getElementsByClassName(cartStyles.productsListElem)) {
            const productID = parseInt(itemContainer.getAttribute('product_id'));
            itemContainer.getElementsByClassName(cartStyles.imageWrapper)[0]
                .addEventListener('click', () => {
                    Bus.globalBus.emit(Events.ProductChangeID, productID);
                    Router.open(`/item/${productID}`);
                });
        }

        document.getElementsByClassName(cartStyles.orderButtonWrapper)[0].addEventListener('click', (evt) => {
            evt.preventDefault();
            Router.open('/order');
        });
    };

    /**
     *
     * @param {number} changedID
     * @param {number} diff
     */
    changeContent = (changedID, diff) => {
        const product = this.presenter.products.find((elem) => {
            return elem.id === changedID;
        });


        document.getElementById('baseCost').innerHTML =
            (parseInt(document.getElementById('baseCost').innerHTML) +
                diff * product.price.base_cost).toString() + '₽';
        document.getElementById('totalCost').innerHTML =
            (parseInt(document.getElementById('totalCost').innerHTML) +
                diff * product.price.total_cost).toString() + '₽';

        document.getElementById('discount').innerHTML = '- ' +
            (parseInt(document.getElementById('baseCost').innerHTML) -
            parseInt(document.getElementById('totalCost').innerHTML)).toString() + '₽';


        const item = Array.from(document.getElementsByClassName(cartStyles.productsListElem)).find((elem) => {
            return +elem.getAttribute('product_id') === changedID;
        });
        const counter = item.getElementsByClassName(cartStyles.count)[0];
        counter.textContent = (+counter.textContent + diff).toString();
        if (!+counter.textContent) {
            item.remove();
            return;
        }

        document.getElementById('totalCost' + changedID.toString()).innerHTML =
            (parseInt(document.getElementById('totalCost' + changedID.toString()).innerHTML) +
                diff * product.price.total_cost).toString() + '₽';
        document.getElementById('baseCost' + changedID.toString()).innerHTML =
            (parseInt(document.getElementById('baseCost' + changedID.toString()).innerHTML) +
                diff * product.price.base_cost).toString() + '₽';
        document.getElementById('discount' + changedID.toString()).innerHTML = 'Скидка ' +
            (parseInt(document.getElementById('baseCost' + changedID.toString()).innerHTML) -
             parseInt(document.getElementById('totalCost' + changedID.toString()).innerHTML)).toString() + '₽';
    }

    /**
     * @param {number | string}id
     */
    removeProduct = (id) => {
        const item = Array.from(document.getElementsByClassName(cartStyles.productsListElem)).filter((item) => {
            return +item.getAttribute('product_id') === +id;
        })[0];
        if (item) {
            item.remove();
        }
    }
}

export default CartView;
