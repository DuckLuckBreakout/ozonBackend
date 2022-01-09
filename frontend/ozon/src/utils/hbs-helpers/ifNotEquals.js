/**
 * @param {Object} arg1
 * @param {Object} arg2
 * @param {Object} options
 * @return {*}
 */
const ifNotEquals = function(arg1, arg2, options) {
    return (arg1 === arg2) ? options.inverse(this) : options.fn(this);
};

export default ifNotEquals;
