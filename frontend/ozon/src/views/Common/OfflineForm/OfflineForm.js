import offlineFormTemplate from './OfflineForm.hbs';
import offlineFormStyles from './OfflineForm.scss';
import imageStyles from '../Img/Img.scss';


/**
 * @class OfflineForm
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class OfflineForm {
    /**
     * A default constructor for OfflineForm
     */
    constructor() {
        this.objectType = 'offlineForm';
    };

    /**
     *
     * @return {string} HTML string, which will generated after templating
     */
    getHtmlString = () => {
        return offlineFormTemplate({
            formType: this.objectType,
            styles: offlineFormStyles,
            imageStyles: imageStyles,
        });
    };
}

export default OfflineForm;
