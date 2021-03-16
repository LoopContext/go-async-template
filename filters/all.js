const _ = require('lodash');

const filter = module.exports;

const typeMap = new Map();
typeMap.set('boolean', 'boolean');
typeMap.set('integer', 'int64');
typeMap.set('null', 'string');
typeMap.set('number', 'float64');
typeMap.set('string', 'string');

filter.checkHasHeaders = function() {
    return false;
}

filter.channelInfo = function([channel, operation]) {

    const channelMeta = {
        funcName: _.upperFirst(operation.id()),
        operation: operation,
        payload: operation.message().payload(),
        payloadStructName: '',
        headers: operation.message().headers(),
        hasHeaders: false,
        headerStructName: _.upperFirst(operation.id() + 'Headers'),
        handlerStructName: _.upperFirst(operation.id() + 'Handler')
    };

    channelMeta.payloadStructName = _.upperFirst(channelMeta.payload.ext('x-parser-schema-id'))

    if (channelMeta.headers) {
        headerProps = channelMeta.headers.properties();
        channelMeta.hasHeaders = Object.keys(headerProps).length > 0;
    }

    return channelMeta;
}

filter.fixPropName = function(propertyName) {
    let usePropName = '';
    _.forEach(propertyName.replace(new RegExp('-', 'g'), ' ').split(' '), propNamePart => {
        usePropName += _.upperFirst(propNamePart);
    }) ;
    return usePropName;
}

filter.fixType = function([name, required, property]) {

    let type = property.type;
    if (typeof type == "function") {
        type = property.type();
    }

    if (!required) {
        required = []
    }

    let typeName = '';

    if (type === 'array') {
        if (!property.items) {
            throw new Error("Array named " + name + " must have an 'items' property to indicate what type the array elements are.")
        }
        const items = property.items();
        if (items.length < 1) {
            throw new Error("Array named " + name + " must have an 'items' property with at least one entry.")
        }
        let item = items[0];
        let itemsType = item.type();
        if (itemsType) {
            if (itemsType === 'object') {
                itemsType = _.upperFirst(item.ext('x-parser-schema-id'));
            } else {
                itemsType = typeMap.get(itemsType);
            }
        }
        typeName = '[]*' + itemsType;
    } else if (type === 'object') {
        typeName = _.upperFirst(property.ext('x-parser-schema-id'));
    } else {
        typeName = typeMap.get(type)
        let format = property.format();
        if (format === 'date') {
            typeName = 'common.Date'
        } else if (format === 'date-time') {
            typeName = 'time.Time'
        }
    }

    if (required.indexOf(name) === -1) {
        typeName = '*' + typeName
    }

    if (!typeName) {
        throw new Error("Can't determine the type of property " + name);
    }

    return typeName
}

filter.checkAttrs = function(asyncapi) {
    const asyncAttrs = {
        needImport:  false,
        hasPublish: false,
        hasSubscribe: false,
        importDate:  false,
        importTime:  false,
    };
    _.forEach(asyncapi.channels(), channel => {
        if (channel.hasPublish()) {
            asyncAttrs.hasPublish = true;
        }
        if (channel.hasSubscribe()) {
            asyncAttrs.hasSubscribe = true;
        }
    });
    _.forEach(asyncapi.components().schemas(), schema => {
        _.forEach(schema.properties(), property => {
            const format = property.format();
            if (format === 'date') {
                asyncAttrs.importDate = true;
            } else if (format === 'date-time') {
                asyncAttrs.importTime = true;
            }
        });
    });
    return asyncAttrs;
}

filter.upperFirst = function (schemaName) {
    return _.upperFirst(schemaName);
}

filter.indent1 = function () {
    return '    ';
}
