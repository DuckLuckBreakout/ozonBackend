import BaseView from '../BaseView.js';
import Events from '../../utils/bus/events';
import orderTemplate from './OrderView.hbs';
import noticeTemplate from './OrderNotice.hbs';
import noticeStyles from './OrderNotice.scss';
import textStyles from './../Common/TextArea/TextArea.scss';
import Input from '../Common/Input/Input';
import inputStyles from '../Common/Input/Input.scss';
import orderStyles from './OrderView.scss';
import buttonStyles from './../Common/Button/Button.scss';
import popupStyles from './../Common/Popup/Popup.scss';
import linkStyles from './../Common/Link/Link.scss';
import Router from '../../utils/router/Router';
import Popup from '../Common/Popup/Popup';
import Blind from '../Common/Blind/Blind';
import decorators from '../decorators.scss';
import Slider from '../Common/Slider/Slider';
import productTemplate from './OrderProduct.hbs';
import imgStyles from '../Common/Img/Img.scss';
import Bus from '../../utils/bus/bus';


/**
 * @class ProductsView
 * @extends BaseView
 * @classdesc Class for showing product
 */
class OrderView extends BaseView {
    /**
     * @param {Object} URLParams
     */
    show = (URLParams = {}) => {
        this.bus.emit(Events.OrderLoad);
    }

    render = () => {
        this.goUp();
        this.parent.innerHTML = '';
        const template = orderTemplate({
            productsList: this.presenter.products,
            inputFields: [new Input({type: 'email', name: 'email', placeholder: 'Электронная почта',
                value: this.presenter.recipient.email}),
            new Input({type: 'text', name: 'firstName', placeholder: 'Имя',
                value: this.presenter.recipient.first_name}),
            new Input({type: 'text', name: 'lastName', placeholder: 'Фамилия',
                value: this.presenter.recipient.last_name}),
            new Input({type: 'text', name: 'address', placeholder: 'Адрес доставки',
                value: this.presenter.address}),
            ],
            orderStyles: orderStyles,
            buttonStyles: buttonStyles,
            textStyles: textStyles,
            linkStyles: linkStyles,
            inputStyles: inputStyles,
            decorators: decorators,
            cartSize: this.presenter.products.length,
            baseCost: +this.presenter.price.total_base_cost,
            discount: +this.presenter.price.total_discount,
            totalCost: +this.presenter.price.total_cost,
        });
        this.cache = new DOMParser().parseFromString(template, 'text/html').getElementById('products-list-block');
        this.parent.appendChild(this.cache);

        const items = [];
        this.presenter.products.forEach((product) => {
            items.push(productTemplate({
                src: product.preview_image,
                id: product.id,
                orderStyles: orderStyles,
                imgStyles: imgStyles,
            }));
        });
        const slider = new Slider(items);
        const recommendationsBlock = document.getElementsByClassName(orderStyles.images)[0];
        recommendationsBlock.appendChild(slider.render());
        slider.checkOverflow();

        for (const itemContainer of this.cache.getElementsByClassName(orderStyles.image)) {
            const productID = parseInt(itemContainer.getAttribute('item-id'));
            itemContainer.addEventListener('click', () => {
                Bus.globalBus.emit(Events.ProductChangeID, productID);
                Router.open(`/item/${productID}`);
            });
        }

        document.getElementById('promo-btn').addEventListener('click', () => {
            this.bus.emit(Events.SendPromo);
        });

        this.cache.getElementsByClassName(buttonStyles.order)[0].addEventListener('click', (evt) => {
            evt.preventDefault();
            this.bus.emit(Events.SendOrder);
            const notice = new DOMParser().parseFromString(new Popup().getHtmlString({
                popupBody: noticeTemplate({
                    styles: noticeStyles,
                    textStyles: textStyles,
                }),
                background: new Blind().getHtmlString(),
                popupType: popupStyles.order,
            }), 'text/html').getElementById('popup');
            const body = document.getElementsByTagName('body')[0];
            body.classList.add(decorators.noScroll);
            this.parent.appendChild(notice);
            document.getElementById('blind')
                .addEventListener('click', (evt) => {
                    evt.preventDefault();
                    body.classList.remove(decorators.noScroll);
                    document.getElementById('popup').remove();
                    Router.open('/orders', {replaceState: true});
                });
        });
    };

    /**
     * @param {Object} newBillInfo
     * @param {String} promocode
     */
    drawNewBill = (newBillInfo, promocode) => {
        const info = document.getElementsByClassName(orderStyles.orderInfo)[0];
        info.getElementsByClassName(textStyles.orderDiscount)[0].innerHTML = `- ${newBillInfo.total_discount} ₽`;
        info.getElementsByClassName(orderStyles.totalPriceBlock)[0].
            getElementsByTagName('span')[1].innerHTML = `${newBillInfo.total_cost} ₽`;
        const promoStatus = document.getElementById('promo-status');
        promoStatus.className = orderStyles.promoStatusSuccess;
        promoStatus.innerHTML = `Промокод ${promocode} был успешно применён!`;
    }

    /**
     * @param {string} promocode
     */
    drawIncorrectPromo = (promocode) => {
        const promoStatus = document.getElementById('promo-status');
        promoStatus.className = orderStyles.promoStatusFailure;
        promoStatus.innerHTML = `Промокода ${promocode} не существует!`;
    }

    /**
     * @param {string} promocode
     */
    drawNoPromoProducts = (promocode) => {
        const promoStatus = document.getElementById('promo-status');
        promoStatus.className = orderStyles.promoStatusFailure;
        promoStatus.innerHTML = `В вашем заказе нет товаров для которых применим промокод ${promocode}!`;
    }
}

export default OrderView;
