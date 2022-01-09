import blindTemplate from './Blind.hbs';
import blindStyles from './Blind.scss';

/**
 * @class Blind
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Blind {
    /**
     * A default constructor for Blind
     */
    constructor() {
        this.objectType = 'blind';
    };

    /**
     * @return {string} HTML string, which will generated after templating
     */
    getHtmlString = () => {
        return blindTemplate({blindStyles});
    };
}

export default Blind;
