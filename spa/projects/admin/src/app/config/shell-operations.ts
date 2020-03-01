import { ShellOperation } from '../shell-operations/model';
import * as ConfigActions from './+state/config.actions';

export const ConfigUpdateShellOperation: ShellOperation = {
    opKey: 'config-update',
    title: 'Update Configuration',
    description: 'Changes in configuration values',
    successAction: ConfigActions.updateConfigSuccess.type,
    failureAction: ConfigActions.updateConfigFailure.type,
};
