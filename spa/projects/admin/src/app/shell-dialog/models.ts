import { Type } from '@angular/core';

export interface ShellDialogConfig {
    title: string;
    component: Type<any>;
    extra?: any;
}
