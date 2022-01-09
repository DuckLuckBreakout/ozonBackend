import popupTemplate from './Popup.hbs';
import popupStyles from './Popup.scss';

/**
 * @class Popup
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Popup {
    /**
     * @description Default constructor
     */
    constructor({} = {}) {
        this.objectType = 'popup';
    };

    /**
     *
     * @param {Object[]} popupBody elements which will appear in popup
     * @param {Object} background background of popup
     * @param {string} popupType type of popup
     * @return {string} generated HTML after templating
     */
    getHtmlString = ({popupBody, background, popupType}) => {
        return popupTemplate({popupBody, background, popupType, popupStyles});
    };
}

export default Popup;
