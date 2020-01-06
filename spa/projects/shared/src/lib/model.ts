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
}

export interface InvoiceCreationParams {
    total: number;
    currency: string;
}
