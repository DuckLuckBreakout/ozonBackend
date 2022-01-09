import Img from '../Img/Img';
import ListOfProductsItem from './ListOfProductsItem/ListOfProductsItem';
import listOFProductsItemTemplate from './ListOfProducts.hbs';
import listOFProductsStyles from './ListOfProducts.scss';

/**
 * @class ListOfProducts
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class ListOfProducts {
    /**
     *
     * @param {Object[]} items array of items to present on page
     * @param {String} type
     */
    constructor(items, type) {
        this.items = [];
        items.forEach((item) => {
            const base = item['price']['base_cost'];
            const discount = item['price']['discount'];
            const discountPrice = Math.ceil(base * (1 - discount * 0.01));
            this.items.push(new ListOfProductsItem({
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
                type: type,
            }));
        });
    };

    /**
     *
     * @return {string} generated HTML after templating
     */
    getHtmlString = () => {
        const itemsTemplates = [];
        this.items.forEach((item) => {
            itemsTemplates.push(item.getHtmlString());
        });

        return listOFProductsItemTemplate({
            items: itemsTemplates,
            listOFProductsStyles: listOFProductsStyles,
        });
    };
}

export default ListOfProducts;
