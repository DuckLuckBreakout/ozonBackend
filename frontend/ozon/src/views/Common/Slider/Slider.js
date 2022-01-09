import sliderTemplate from './Slider.hbs';
import sliderStyles from './Slider.scss';
import buttonStyles from './../Button/Button.scss';
import decorators from '../../decorators.scss';
import Img from '../Img/Img';

/**
 * @class Slider
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Slider {
    /**
     *
     * @param {Object[]} items
     */
    constructor(items) {
        this.items = items;
        this.content = null;
        this.contentGap = null;
        this.contentWrapper = null;
        this.left = null;
        this.right = null;
        this.itemWidth = null;
        this.currentTranslate = 0;
    };

    /**
     *
     * @return {HTMLElement}
     */
    render = () => {
        const template = sliderTemplate({
            items: this.items,
            sliderStyles: sliderStyles,
            buttonStyles: buttonStyles,
            decorators: decorators,
            buttonSvg: new Img({src: '/svg/slider.svg'}),
        });
        const slider = new DOMParser().parseFromString(template, 'text/html');
        this.content = slider.getElementsByClassName(sliderStyles.content)[0];
        this.contentWrapper = slider.getElementsByClassName(sliderStyles.contentWrapper)[0];
        this.right = slider.getElementById('slider-right');
        this.left = slider.getElementById('slider-left');

        Array.from(slider.getElementsByClassName(buttonStyles.slider)).forEach((button) => {
            button.addEventListener('click', () => {
                button.disabled = true;
                const buttonTimeout = parseFloat(
                    getComputedStyle(this.content)
                        .transitionDuration.match(/[0-9]+\.?[0-9]+/)[0],
                ) * 1000;
                setTimeout(() => {
                    button.disabled = false;
                }, buttonTimeout);

                const contentWrapperBox = this.contentWrapper.getBoundingClientRect();
                const contentBox = this.content.getBoundingClientRect();

                if (button === this.left) {
                    this.currentTranslate += Math.floor(contentWrapperBox.width / this.itemWidth) * this.itemWidth;
                } else {
                    this.currentTranslate -= Math.floor(contentWrapperBox.width / this.itemWidth) * this.itemWidth;
                }

                if (this.currentTranslate >= 0) {
                    this.currentTranslate = 0;
                    this.left.classList.add(decorators.noVisibility);
                } else {
                    this.left.classList.remove(decorators.noVisibility);
                }

                if (-this.currentTranslate >= this.content.scrollWidth - contentBox.width) {
                    this.currentTranslate = contentBox.width - this.content.scrollWidth;
                    this.right.classList.add(decorators.noVisibility);
                } else {
                    this.right.classList.remove(decorators.noVisibility);
                }

                this.content.style.transform = `translateX(${this.currentTranslate}px)`;
            });
        });
        return slider.getElementById('slider');
    }

    checkOverflow = () => {
        this.contentGap = parseFloat(getComputedStyle(this.content).gap.match(/[0-9]+\.?[0-9]+/)[0]);
        this.itemWidth = this.content.children[0].getBoundingClientRect().width + this.contentGap;
        if (this.contentWrapper.scrollWidth > this.content.getBoundingClientRect().width) {
            this.content.style.maxWidth = `${Math.floor(
                this.content.getBoundingClientRect().width / this.itemWidth,
            ) * this.itemWidth}px`;

            this.right.classList.remove(decorators.noVisibility);
        }
    }
}

export default Slider;
