import * as uuid_ from 'uuid/v4';
// https://github.com/jvandemo/generator-angular2-library/issues/221#issuecomment-355945207
// :manshrug
const uuid = uuid_;

export function generateUniqueId() {
    return uuid();
}
