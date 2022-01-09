import Button from '../Button/Button.js';
import paginationTemplate from './Pagination.hbs';
import paginationOneButtonTemplate from './PaginationOneButton.hbs';
import styles from './Pagination.scss';
import buttonStyles from '../Button/Button.scss';
import textStyles from '../TextArea/TextArea.scss';

/**
 * @class Pagination
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Pagination {
    /**
     * @param {number} pagesCount maximum amount of page
     * @param {boolean} isOneButton
     * @param {number} currentPage current page!
     */
    constructor({pagesCount, currentPage}, isOneButton = false) {
        this.pagesCount = pagesCount;
        this.currentPage = parseInt(currentPage);
        this.maxButtons = 10;
        this.isOneButton = isOneButton;
    }

    /**
     *
     * @return {string} generated HTML after templating
     */
    getHtmlString = () => {
        if (this.isOneButton) {
            return this.getHtmlOneButtonString();
        }

        if (!this.pagesCount) {
            return '';
        }

        const buttons = [];
        let surplusPages = 0;
        let shortagePages = 0;
        for (
            let i = this.currentPage - this.maxButtons / 2;
            i <= this.currentPage + this.maxButtons / 2 + surplusPages;
            i++
        ) {
            if (i <= 0) {
                surplusPages++;
                continue;
            }
            if (i > this.pagesCount) {
                shortagePages = this.currentPage + this.maxButtons / 2 + surplusPages - i;
                break;
            }
            const buttonName = i === this.currentPage ? buttonStyles.paginatorCurrent : buttonStyles.paginatorCommon;
            buttons.push(new Button({
                name: buttonName,
                value: i,
            }));
        }

        const shortageButtons = [];
        for (
            let i = parseInt(buttons[0].value) - shortagePages;
            i <= parseInt(buttons[0].value) + shortagePages;
            i++
        ) {
            if (i <= 0) {
                continue;
            }
            if (parseInt(buttons[0].value) <= i) {
                break;
            }
            shortageButtons.push(new Button({
                name: styles.commonPage,
                value: i,
            }));
        }
        const resultButtons = shortageButtons.concat(buttons);
        const firstButton = new Button({
            name: styles.bigButton,
            value: 1,
        });
        const nextButton = new Button({
            name: styles.bigButton,
            value: this.currentPage + 1,
        });
        if (parseInt(nextButton.value) > this.pagesCount) {
            nextButton.value = -1;
        }

        return paginationTemplate({
            firstButton: firstButton,
            buttons: resultButtons,
            nextButton: nextButton,
            styles: styles,
            buttonStyles: buttonStyles,
            textStyles: textStyles,
        });
    }
    /**
     * @return {string} HTML
     */
    getHtmlOneButtonString() {
        if (this.currentPage === this.pagesCount || !this.pagesCount) {
            return '';
        }

        const loadMoreButton = new Button({
            name: buttonStyles.paginatorBig,
            value: 'Загрузить ещё',
        });
        return paginationOneButtonTemplate({
            buttons: [loadMoreButton],
            buttonStyles: buttonStyles,
            styles: styles,
        });
    }
}

export default Pagination;
