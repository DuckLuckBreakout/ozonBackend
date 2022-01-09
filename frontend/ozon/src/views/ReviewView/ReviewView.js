import BaseView from '../BaseView.js';
import Events from '../../utils/bus/events';
import reviewTemplate from './ReviewView.hbs';
import Input from '../Common/Input/Input';
import reviewStyles from './ReviewView.scss';
import textStyles from './../Common/TextArea/TextArea.scss';
import imgStyles from './../Common/Img/Img.scss';
import buttonStyles from './../Common/Button/Button.scss';
import linkStyles from './../Common/Link/Link.scss';
import Router from '../../utils/router/Router';
import Bus from '../../utils/bus/bus';
import Popup from '../Common/Popup/Popup';
import noticeTemplate from './ReviewNotice.hbs';
import noticeStyles from './ReviewNotice.scss';
import Blind from '../Common/Blind/Blind';
import popupStyles from '../Common/Popup/Popup.scss';
import decorators from '../decorators.scss';
import Img from '../Common/Img/Img';

/**
 * @class ReviewView
 * @extends BaseView
 * @classdesc Class for showing review
 */
class ReviewView extends BaseView {
    show = () => {
        this.presenter.rating = 0;
        this.presenter.isPublic = true;
        this.render();
    }

    render = () => {
        this.parent.innerHTML = '';
        const userName = this.presenter.userName === '' ||
                         this.presenter.userName === ' ' ||
                         this.presenter.userName === 'null null' ||
                         this.presenter.userName === undefined ?
            'Анонимный пользователь' :
            this.presenter.userName;
        const template = reviewTemplate({
            inputFields: [
                new Input({
                    type: 'text',
                    name: 'advantages',
                    placeholder: 'Достоинства',
                }),
                new Input({
                    type: 'text',
                    name: 'disadvantages',
                    placeholder: 'Недостатки',
                }),
                new Input({
                    type: 'text',
                    name: 'comment',
                    placeholder: 'Комментарий к отзыву',
                }),
            ],
            rating: 3,
            isPublic: new Input({
                type: 'checkbox',
                name: 'isPublic',
            }),
            styles: reviewStyles,
            textStyles: textStyles,
            buttonStyles: buttonStyles,
            imgStyles: imgStyles,
            linkStyles: linkStyles,
            product: this.presenter.product,
            userName: userName,
            ratingStar: new Img({src: '/svg/empty_star.svg'}),
        });
        this.cache = new DOMParser().parseFromString(template, 'text/html').getElementById('review-block');
        this.parent.appendChild(this.cache);

        const emptyStarLink = '/svg/empty_star.svg';
        const starLink = '/svg/star.svg';
        const stars = Array.from(this.cache.getElementsByClassName(reviewStyles.stars)[0].children);

        const emptyRating = (event) => {
            if (!event.target.classList.contains(reviewStyles.stars) &&
                !event.target.classList.contains(reviewStyles.star)) {
                stars.forEach((star, i) =>
                    star.setAttribute('src', i < this.presenter.rating ? starLink : emptyStarLink),
                );
            }
        };

        this.cache.addEventListener('mouseover', emptyRating);


        stars.forEach((star) => {
            const currentStarValue = parseInt(star.getAttribute('value'));
            star.addEventListener('mouseover', () => {
                for (let i = 0; i < 5; i++) {
                    const isStar = parseInt(stars[i].getAttribute('value')) <= currentStarValue;
                    stars[i].setAttribute('src', isStar ? starLink : emptyStarLink);
                }
            });
            star.addEventListener('click', () => {
                this.presenter.rating = currentStarValue + 1;
            });
        });

        const submitButton = this.cache.getElementsByClassName(buttonStyles.review)[0];
        submitButton.addEventListener('click', (event) => {
            event.preventDefault();
            this.presenter.isPublic = !document.getElementsByName('isPublic')[0].checked;
            this.bus.emit(Events.ReviewOrder);
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
                    Router.open(`/item/${this.presenter.product.id}`, {replaceState: true});
                });
        });

        const productLink = this.cache.getElementsByClassName(linkStyles.link)[0];
        productLink.addEventListener('click', () => {
            const productId = this.presenter.product.id;
            Bus.globalBus.emit(Events.ProductChangeID, productId);
            Router.open(`/item/${productId}`);
        });

        if (this.presenter.userName === '' ||
            this.presenter.userName === ' ' ||
            this.presenter.userName === 'null null' ||
            this.presenter.userName === undefined) {
            document.getElementsByName('isPublic')[0].checked = true;
            document.getElementsByClassName(reviewStyles.private)[0].style.display = 'none';
        }
    };
}

export default ReviewView;
