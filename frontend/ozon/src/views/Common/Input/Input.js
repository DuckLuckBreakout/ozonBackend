/**
 * @class Input
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Input {
    /**
     *
     * @param {string} type type of input (email, password, etc...)
     * @param {string} name name of an input in html tag
     * @param {string} placeholder placeholder of an input in html tag
     * @param {string} value value of an input in html tag
     * @param {boolean} isDisabled if true - input field disabled. Otherwise - false.
     */
    constructor({type = 'text',
        name = '',
        placeholder='',
        value = '',
        isDisabled=false} = {}) {
        this.objectType = 'input';
        this.name=name;
        this.type = type;
        this.placeholder = placeholder;
        this.value = value;
        this.isDisabled = isDisabled;
    }
}

export default Input;
