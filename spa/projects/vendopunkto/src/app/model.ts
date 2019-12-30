export interface AtomicUnit {
    value: number;
    valueFormatted: string;
}

export interface InvoiceDTO {
    status: number;
    id: string;
    total: AtomicUnit;
    remaining: AtomicUnit;
    currency: string;
    createdAt: string;
    paymentMethods?: PaymentMethod[];
    paymentPercentage: number;
    payments?: Payment[];
}

export interface PaymentMethod {
    id: number;
    total: AtomicUnit;
    currency: string;
    address: string;
    remaining: AtomicUnit;
    qrCode: string;
}

export interface Payment {
    status: number;
    txHash: string;
    amount: AtomicUnit;
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
