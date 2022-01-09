import BaseView from '../BaseView.js';
import Img from '../Common/Img/Img.js';
import productPageTemplate from './ProductView.hbs';
import reviewsTemplate from './ProductReviews.hbs';
import Events from '../../utils/bus/events';
import productStyles from './ProductView.scss';
import reviewStyles from './ProductReview.scss';
import buttonStyles from './../Common/Button/Button.scss';
import textStyles from './../Common/TextArea/TextArea.scss';
import imgStyles from './../Common/Img/Img.scss';
import linkStyles from './../Common/Link/Link.scss';
import decorators from '../decorators.scss';
import Bus from '../../utils/bus/bus';
import Router from '../../utils/router/Router';
import Pagination from '../Common/Pagination/Pagination';
import starsTemplate from './ProductStars.hbs';
import starsStyles from './ProductStars.scss';
import Slider from '../Common/Slider/Slider';
import ListOfProductsItem from '../Common/ListOfProducts/ListOfProductsItem/ListOfProductsItem';
import ListOfProductsItemStyles from '../Common/ListOfProducts/ListOfProductsItem/ListOfProductsItem.scss';

/**
 * @class ProductView
 * @extends BaseView
 * @classdesc Class for Product page
 */
class ProductView extends BaseView {
    /**
     * @param {Object} URLParams
     */
    show = (URLParams = {}) => {
        if (!this.IDs) {
            this.IDs = {};
        }
        if (!this.IDs['productID']) {
            this.IDs['productID'] = 1;
        }

        this.presenter.sortKey = 'date';
        this.presenter.sortDirection = 'DESC';
        this.bus.emit(Events.ProductLoad, this.IDs['productID']);
    }

    /**
     *
     */
    render = () => {
        this.goUp();
        this.parent.innerHTML = '';
        const images = [];
        this.presenter.item['images'].forEach((src) => {
            images.push(new Img({
                src: src,
            }));
        });

        const template = productPageTemplate({
            name: this.presenter.item['name'],
            price: this.presenter.item['price'],
            rating: this.presenter.item['rating'],
            width: this.presenter.item['rating'] / 5 * 100,
            images: images,
            description: this.presenter.item['description']['descriptionText'],
            category: this.presenter.item['description']['category'],
            productStyles: productStyles,
            reviewStyles: reviewStyles,
            buttonStyles: buttonStyles,
            textStyles: textStyles,
            imgStyles: imgStyles,
            linkStyles: linkStyles,
            decorators: decorators,
            category_path: this.presenter.item.category_path.reverse(),
            properties: this.presenter.item['properties'],
            propertiesMain: this.presenter.item['properties']['Общие'],
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
        });
        this.cache = new DOMParser().parseFromString(template, 'text/html')
            .getElementsByClassName(productStyles.block)[0];
        this.parent.appendChild(this.cache);

        const descriptionBlock = this.cache.getElementsByClassName(productStyles.description)[0];
        if (this.presenter.item['description']['descriptionText'] === '') {
            descriptionBlock.className = '';
            descriptionBlock.classList.add(decorators.hidden);
        }

        const propertiesBlock = this.cache.getElementsByClassName(productStyles.properties)[0];
        if (Object.keys(this.presenter.item['properties']).length === 0) {
            propertiesBlock.className = '';
            propertiesBlock.classList.add(decorators.hidden);
        }

        const select = document.getElementsByTagName('select')[0];
        select.addEventListener('change', () => {
            const selected = select.selectedOptions[0];
            const sortKey = selected.getAttribute('key');
            const sortDirection = selected.getAttribute('direction');
            this.presenter.sortKey = sortKey;
            this.presenter.sortDirection = sortDirection;
            this.cache.getElementsByClassName(productStyles.reviewList)[0].innerHTML = '';
            Bus.globalBus.emit(Events.GetProductReviews, +this.IDs['productID'], 1,
                this.presenter.sortKey,
                this.presenter.sortDirection);
        });

        const mainImage = this.cache.getElementsByClassName(imgStyles.productMain)[0];
        Array.from(this.cache.getElementsByClassName(imgStyles.productCommon)).forEach((image) => {
            image.addEventListener('click', () => {
                mainImage.setAttribute('src', image.getAttribute('src'));
            });
        });

        Array.from(this.cache.getElementsByClassName(linkStyles.link)).forEach((category) => {
            category.addEventListener('click', () => {
                Router.open(`/items/${category.getAttribute('category_id')}`);
            });
        });


        const button = document.getElementsByClassName(buttonStyles.notInCartProduct)[0];
        button.addEventListener('click', () => {
            Bus.globalBus.emit(Events.CartAddProduct, this.IDs['productID'], 1, Events.ProductItemAdded);
        });

        const reviewButton = this.cache.getElementsByClassName(buttonStyles.review)[0];
        if (reviewButton) {
            reviewButton.addEventListener('click', () => {
                Bus.globalBus.emit(Events.ChangeReviewProductId, this.IDs['productID']);
                Router.open('/review');
            });
        }

        Bus.globalBus.emit(Events.CartGetProductID);
        Bus.globalBus.emit(Events.GetProductReviews, +this.IDs['productID'], 1,
            this.presenter.sortKey,
            this.presenter.sortDirection);

        Bus.globalBus.emit(Events.ChangeReviewProductId, this.IDs['productID']);
        Bus.globalBus.emit(Events.ReviewRightsLoad);
        this.bus.emit(Events.ProductStarsCounterLoad);
    }


    setProductAdded = () => {
        if (!this.cache) {
            return;
        }
        const button = this.cache.getElementsByClassName(buttonStyles.notInCartProduct)[0];
        if (!button) {
            return;
        }
        const newButton = button.cloneNode(true);
        button.replaceWith(newButton);

        newButton.className = buttonStyles.inCartProduct;
        newButton.getElementsByTagName('span')[0].innerHTML = 'В корзине';
        newButton.addEventListener('click', () => {
            Bus.globalBus.emit(Events.CartRemoveProduct, this.IDs['productID'], Events.ProductItemNotAdded);
        });
    }

    setProductNotAdded = () => {
        if (!this.cache) {
            return;
        }
        const button = document.getElementsByClassName(buttonStyles.inCartProduct)[0];
        if (!button) {
            return;
        }
        const newButton = button.cloneNode(true);
        button.replaceWith(newButton);

        newButton.getElementsByTagName('span')[0].innerHTML = 'В корзину';
        newButton.className = buttonStyles.notInCartProduct;
        newButton.addEventListener('click', () => {
            Bus.globalBus.emit(Events.CartAddProduct, this.IDs['productID'], 1, Events.ProductItemAdded);
        });
    }

    /**
     * @param {Array} reviews
     * @param {Object} paginationInfo
     */
    renderProductsReview = (reviews, paginationInfo) => {
        if (reviews.length === 0 && paginationInfo.pagesCount === 0) {
            const reviewsWrapper = this.cache.getElementsByClassName(productStyles.reviewsWrapper)[0];
            reviewsWrapper.className = '';
            reviewsWrapper.classList.add(productStyles.zeroReviews);
            reviewsWrapper.innerHTML = 'Пользователи не оставили отзывов для данного товара';
        }

        const pagination = new Pagination(paginationInfo, true).getHtmlString();
        if (document.getElementById('review-pagination') !== undefined &&
            document.getElementById('review-pagination') !== null) {
            document.getElementById('review-pagination').innerHTML = pagination;
        }
        reviews.forEach((review) => {
            review.date_added = review.date_added.slice(0, 10);
            review.user_name = (review.user_name !== '') ? review.user_name : 'Анонимный пользователь';
            review.firstChar = review.user_name[0];
            if (review.user_avatar) {
                review.avatar = `/${review.user_avatar}`;
            }
        });

        if (this.cache.getElementsByClassName(reviewStyles.block)[0]) {
            this.cache.getElementsByClassName(reviewStyles.block)[0].innerHTML +=
                reviewsTemplate({
                    reviewsList: reviews,
                    productStyles: productStyles,
                    reviewStyles: reviewStyles,
                    imgStyles: imgStyles,
                    emptyStarSvg: new Img({src: '/svg/empty_star.svg'}),
                    starSvg: new Img({src: '/svg/star.svg'}),
                });
        } else {
            this.cache.getElementsByClassName(productStyles.reviewList)[0].innerHTML = reviewsTemplate({
                reviewsList: reviews,
                pagination: pagination,
                productStyles: productStyles,
                reviewStyles: reviewStyles,
                imgStyles: imgStyles,
                emptyStarSvg: new Img({src: '/svg/empty_star.svg'}),
                starSvg: new Img({src: '/svg/star.svg'}),
            });
        }
        if (paginationInfo.pagesCount === paginationInfo.currentPage) {
            return;
        }

        const buttonPagination = document.getElementById('review-pagination').getElementsByTagName('button')[0];
        if (buttonPagination) {
            buttonPagination.addEventListener('click', () => {
                Bus.globalBus.emit(Events.GetProductReviews, +this.IDs['productID'], paginationInfo.currentPage + 1,
                    this.presenter.sortKey,
                    this.presenter.sortDirection);
            });
        }
    }

    addReviewButton = () => {
        const reviewButton = document.createElement('button');
        reviewButton.className = buttonStyles.review;
        const reviewSpan = document.createElement('span');
        reviewSpan.className = textStyles.bigButton;
        reviewSpan.innerHTML = 'Оставить отзыв';
        reviewButton.appendChild(reviewSpan);
        reviewButton.addEventListener('click', () => {
            Router.open('/review');
        });
        document.getElementsByClassName(productStyles.reviewButtonWrapper)[0].appendChild(reviewButton);
    }

    renderRecommendations = () => {
        if (this.presenter.recommendations.length === 0) {
            const recommendationsBlock = this.cache.getElementsByClassName(productStyles.recommendationsWrapper)[0];
            recommendationsBlock.className = '';
            recommendationsBlock.classList.add(decorators.hidden);
            return;
        }

        const items = [];
        this.presenter.recommendations.forEach((item) => {
            const base = item['price']['base_cost'];
            const discount = item['price']['discount'];
            const discountPrice = Math.ceil(base * (1 - discount * 0.01));
            items.push(new ListOfProductsItem({
                itemInCart: item['inCart'],
                itemReviewsCount: item['count_reviews'],
                itemId: item['id'],
                itemImage: new Img({src: item['preview_image']}),
                itemName: item['title'],
                itemRating: item['rating'],
                itemPrice: {
                    base: base,
                    discount: discount,
                    discountPrice: discountPrice,
                },
                type: 'recommendations',
            }).getHtmlString());
        });

        const slider = new Slider(items);
        const recommendationsBlock = document.getElementById('recommendations');
        recommendationsBlock.appendChild(slider.render());
        for (const itemContainer of document.getElementsByClassName(ListOfProductsItemStyles.block)) {
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
        Bus.globalBus.emit(Events.CartGetProductsID, Events.ProductCartLoadedProductsID);
        slider.checkOverflow();
    }

    renderStarsCounter = () => {
        const starsInfo = [];
        const rowTitles = ['1 звезда', '2 звезды', '3 звезды', '4 звезды', '5 звезд'];
        const starsCounts = this.presenter.item.stars;
        let starsCountsSum = starsCounts.reduce((accumulator, currentStarCount) => {
            return accumulator + currentStarCount;
        });
        if (starsCountsSum === 0) {
            starsCountsSum = Number.MAX_SAFE_INTEGER;
        }
        for (let i = 0; i < starsCounts.length; i++) {
            starsInfo.push({
                title: rowTitles[i],
                count: starsCounts[i],
                width: starsCounts[i] / starsCountsSum * 100,
            });
        }
        const starsBlock = document.getElementById('rating-stars');
        starsBlock.innerHTML = starsTemplate({
            rating: this.presenter.item.rating.toFixed(2),
            width: this.presenter.item.rating.toFixed(2) / 5 * 100,
            stars: starsInfo.reverse(),
            starsStyles: starsStyles,
            textStyles: textStyles,
            imgStyles: imgStyles,
        });
    }

    /**
     * @param {Set} productsInCart
     */
    recommendationSetAddedProducts = (productsInCart) => {
        for (const item of document.getElementsByClassName(ListOfProductsItemStyles.block)) {
            if (productsInCart.has(+item.getAttribute('item-id'))) {
                this.setButtonAdded(item.getElementsByTagName('button')[0], +item.getAttribute('item-id'));
            }
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
}

export default ProductView;
