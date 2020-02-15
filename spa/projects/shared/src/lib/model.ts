export interface InvoiceDTO {
    status: number;
    id: string;
    total: string;
    remaining: string;
    currency: string;
    createdAt: string;
    paymentMethods?: PaymentMethod[];
    paymentPercentage: number;
    payments?: Payment[];
}

export interface PaymentMethod {
    id: number;
    total: string;
    currency: string;
    address: string;
    remaining: string;
    qrCode: string;
}

export interface Payment {
    status: number;
    txHash: string;
    amount: string;
    confirmations: number;
    blockHeight: number;
    currency: string;
    confirmedAt: Date;
    createdAt: Date;
}

export interface SupportedCurrency {
    name: string;
    symbol: string;
    logoImageUrl: string;
    supportsPayments: boolean;
    supportsPricing: boolean;
}

export interface InvoiceCreationParams {
    total: number;
    currency: string;
    paymentMethods: PaymentMethodCreationParams[];
}

export interface PaymentMethodCreationParams {
    total: number;
    currency: string;
}

export interface GetCurrencyRatesParams {
    fromCurrency: string;
    toCurrencies: string[];
}

export interface GetCurrencyRatesResult {
    [currency: string]: string;
}

export interface GetCurrencyExchangeParams {
    fromCurrency: string;
    toCurrencies: string[];
    amount: string | number;
}

export interface GetCurrencyExchangeResult {
    [currency: string]: string;
}

export interface GetConfigResult {
    currency_metadata_plugin: string;
    default_invoice_currency: string;
    exchange_rates_plugin: string;
    invoice_currencies: string[];
    plugin_hosts: string[];
    wallet_poller_interval: string;
    [k: string]: unknown;
}

export interface GetPluginResult {
    name: string;
    id: string;
    pluginType: string;
}
