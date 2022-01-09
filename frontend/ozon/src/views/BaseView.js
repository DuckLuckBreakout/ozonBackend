/**
 * @class BaseView
 * @classdesc Base class for all other views
 */
class BaseView {
    #parent
    #cache
    #presenter
    #bus
    #ids

    /**
     *
     * @param {Object} parent  Parent HTML element
     * @param {Object} presenter Presenter of this view
     * @param {Object} bus bus of this mvp part
     */
    constructor(parent, presenter, bus) {
        this.#parent = parent;
        this.#parent.dataset.view = this.constructor.name;
        this.#parent.hidden = false;
        this.#cache = '';
        this.#bus = bus;
        this.#ids = {};
        this.#presenter = presenter;
    }

    /**
     *
     * @return {boolean} is element Active?
     */
    get isActive() {
        return !this.parent.hidden;
    }

    /**
     *
     * @return {HTMLElement} Cache of View
     */
    get cache() {
        return this.#cache;
    }

    /**
     *
     * @return {Object} Presenter of View
     */
    get presenter() {
        return this.#presenter;
    }

    /**
     *
     * @return {HTMLElement} Parent of this View
     */
    get parent() {
        return this.#parent;
    }

    /**
     *
     * @return {Object} Bus of this mvp part
     */
    get bus() {
        return this.#bus;
    }

    /**
     *
     * @return {number} id,page number,item id, etc. for this view (use if needs)
     */
    get IDs() {
        return this.#ids;
    }

    /**
     *
     * @param {HTMLElement} cache Cache of View
     */
    set cache(cache) {
        this.#cache = cache;
    }

    /**
     *
     * @param {Object} ids new ids
     */
    set IDs(ids) {
        this.#ids = ids;
    }


    /**
     * @description Hides an element in HTML
     */
    hide = () => {
        this.#cache.hidden = true;
    }

    /**
     * @description Removes Element from HTML
     */
    remove = () => {
        this.#parent.removeChild(this.cache);
    }

    /**
     * @param {Object} URLParams
     * @description Show element in HTML
     */
    show = (URLParams = {}) => {
        this.render();
        this.#cache.hidden = false;
    }

    /**
     * @description clears cache of View
     */
    clearCache = () => {
        this.#cache = '';
    }

    /**
     * @param {HTMLElement} input
     * @param {string} placeholder
     */
    static setInvalidInputPlaceholder(input, placeholder) {
        input.value = '';
        input.placeholder = placeholder;
        input.style['border-color'] = '#ff726f';
    }

    /**
     * @param {array} invalidFields
     */
    invalidForm = (invalidFields) => {
        const inputs = Array.from(this.cache.getElementsByTagName('input'));
        for (const field of invalidFields) {
            const input = inputs.find((input) => {
                return field.fieldName === input.name;
            });
            if (!input) {
                break;
            }
            BaseView.setInvalidInputPlaceholder(input, field.errorMessage);
        }
    }

    /**
     * @description Drop cache and render
     */
    reRender = () => {
        this.clearCache();
        this.render();
    }

    /**
     * @description Raising the user to the top
     */
    goUp = () => {
        window.scrollTo(0, 0);
    }

    /**
     * Render from template
     */
    render = () => {

    }
}

export default BaseView;
