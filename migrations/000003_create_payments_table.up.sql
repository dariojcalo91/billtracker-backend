CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_id UUID NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    month TEXT NOT NULL,
    amount_paid NUMERIC(12,2) NOT NULL,
    proof_file_url TEXT,
    paid_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (bill_id, month)
);

CREATE INDEX idx_payments_bill_id ON payments(bill_id);
CREATE INDEX idx_payments_month ON payments(month);
