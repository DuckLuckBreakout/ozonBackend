import BaseView from '../BaseView.js';
import Events from '../../utils/bus/events';
import ordersTemplate from './OrdersView.hbs';
import ordersListTemplate from './OrdersList.hbs';
import orderImagesTemplate from './OrderImages.hbs';
import ordersStyles from './OrdersView.scss';
import Bus from '../../utils/bus/bus';
import Router from '../../utils/router/Router';
import imgStyles from '../Common/Img/Img.scss';
import textStyles from '../Common/TextArea/TextArea.scss';
import Slider from '../Common/Slider/Slider';
import buttonStyles from '../Common/Button/Button.scss';
import decorators from '../decorators.scss';


/**
 * @class OrdersView
 * @extends BaseView
 * @classdesc Class for showing orderS
 */
class OrdersView extends BaseView {
    #ticking;
    #tickingLock;

    show = () => {
        if (!this.IDs) {
            this.IDs = {};
        }
        if (!this.IDs['page']) {
            this.IDs['page'] = 1;
        }

        if (!this.presenter.sortKey) {
            this.presenter.sortKey = 'date';
        }
        if (!this.presenter.sortDirection) {
            this.presenter.sortDirection = 'DESC';
        }

        this.bus.emit(Events.OrdersLoad,
            this.IDs['page'],
            this.presenter.sortKey,
            this.presenter.sortDirection);
    }

    render = () => {
        this.parent.innerHTML = '';
        const orders = this.presenter.orders;
        const ordersListHtmlString = ordersListTemplate({
            orderList: orders,
            ordersStyles: ordersStyles,
            imgStyles: imgStyles,
            textStyles: textStyles,
        });
        const template = ordersTemplate({
            ordersListHtmlString: ordersListHtmlString,
            isEmpty: !orders.length,
            ordersStyles: ordersStyles,
            select: [
                {
                    key: 'date',
                    direction: 'DESC',
                    name: 'Сначала новые',
                },
                {
                    key: 'date',
                    direction: 'ASC',
                    name: 'Сначала старые',
                },
            ],
            sort: {
                key: this.presenter.sortKey,
                direction: this.presenter.sortDirection,
            },
            imgStyles: imgStyles,
            textStyles: textStyles,
        });
        this.cache = new DOMParser().parseFromString(template, 'text/html').getElementById('orders-list-block');
        this.parent.appendChild(this.cache);
        const selectBlock = this.cache.getElementsByClassName(ordersStyles.select)[0];
        if (this.presenter.orders.length === 0) {
            selectBlock.className = '';
            selectBlock.classList.add(decorators.hidden);
        }

        const images = [];
        orders.forEach((order) => {
            const orderImages = [];
            order.product_images.forEach((img) => {
                orderImages.push(orderImagesTemplate({
                    src: img.preview_image,
                    id: img.id,
                    ordersStyles: ordersStyles,
                    imgStyles: imgStyles,
                }));
            });
            images.push(orderImages);
        });

        const imagesBlocks = this.cache.getElementsByClassName(ordersStyles.images);
        for (let i = 0; i < imagesBlocks.length; i++) {
            const slider = new Slider(images[i]);
            imagesBlocks[i].appendChild(slider.render());
            slider.checkOverflow();
        }


        const select = this.cache.getElementsByClassName(ordersStyles.select)[0];
        select.addEventListener('change', () => {
            const selected = select.selectedOptions[0];
            this.presenter.sortKey = selected.getAttribute('key');
            this.presenter.sortDirection = selected.getAttribute('direction');
            Router.open('/orders');
        });

        for (const itemContainer of this.cache.getElementsByClassName(ordersStyles.image)) {
            const productID = parseInt(itemContainer.getAttribute('item-id'));
            itemContainer.addEventListener('click', () => {
                Bus.globalBus.emit(Events.ProductChangeID, productID);
                Router.open(`/item/${productID}`);
            });
        }

        for (const button of this.cache.getElementsByClassName(buttonStyles.paginator)) {
            button.addEventListener('click', () => {
                const page = parseInt(button.getAttribute('page'));
                this.ID = page;
                Router.open(`/orders/${page}`, {id: page});
            });
        }

        this.#ticking = false;
        window.addEventListener('scroll', (e) => {
            if (+this.presenter.paginationInfo.pagesCount === +this.presenter.paginationInfo.currentPage) {
                return;
            }
            // Scroll Optimisation
            if (!this.#ticking && !this.#tickingLock) {
                window.requestAnimationFrame(() => {
                    const y = window.pageYOffset + window.innerHeight;
                    const scrollHeight = Math.max(
                        document.body.scrollHeight, document.documentElement.scrollHeight,
                        document.body.offsetHeight, document.documentElement.offsetHeight,
                        document.body.clientHeight, document.documentElement.clientHeight,
                    );
                    if (y / scrollHeight > 0.99 && !this.#tickingLock) {
                        this.#tickingLock = true;
                        this.bus.emit(Events.OrdersLoadMoreOrders,
                            this.presenter.sortKey,
                            this.presenter.sortDirection);
                    }
                    this.#ticking = false;
                });
                this.#ticking = true;
            }
        });
    }

    renderMoreOrders = () => {
        const moreOrders = this.presenter.orders;
        const moreOrdersListTemplate = ordersListTemplate({
            orderList: moreOrders,
            ordersStyles: ordersStyles,
            imgStyles: imgStyles,
            textStyles: textStyles,
        });
        const imagesOffset = this.cache.getElementsByClassName(ordersStyles.images).length;
        const childNodes = new DOMParser().
            parseFromString(moreOrdersListTemplate, 'text/html').
            getElementById('new-orders-list');
        document.getElementsByClassName(ordersStyles.orders)[0].appendChild(childNodes);


        const images = [];
        moreOrders.forEach((order) => {
            const orderImages = [];
            order.product_images.forEach((img) => {
                orderImages.push(orderImagesTemplate({
                    src: img.preview_image,
                    id: img.id,
                    ordersStyles: ordersStyles,
                    imgStyles: imgStyles,
                }));
            });
            images.push(orderImages);
        });

        const imagesBlocks = this.cache.getElementsByClassName(ordersStyles.images);
        for (let i = imagesOffset; i < imagesBlocks.length; i++) {
            const slider = new Slider(images[i - imagesOffset]);
            imagesBlocks[i].appendChild(slider.render());
            slider.checkOverflow();
        }

        for (const itemContainer of childNodes.getElementsByClassName(ordersStyles.image)) {
            const productID = parseInt(itemContainer.getAttribute('item-id'));
            itemContainer.addEventListener('click', () => {
                Bus.globalBus.emit(Events.ProductChangeID, productID);
                Router.open(`/item/${productID}`);
            });
        }

        this.#tickingLock = false;
    }
}

export default OrdersView;
