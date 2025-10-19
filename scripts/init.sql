-- Create enum type for message status
CREATE TYPE message_status AS ENUM ('pending', 'sent', 'failed');

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) <= 160),
    status message_status DEFAULT 'pending',
    sent_at TIMESTAMPTZ
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);
CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages(sent_at);

INSERT INTO messages (phone_number, content)
VALUES
('+84901234567', 'Hello from Insider!'),
('+84976543210', 'Your code is 1234'),
('+84968889999', 'Reminder: meeting at 3PM'),
('+84913334567', 'Reminder: meeting at 4PM'),
('+84987777777', 'Reminder: meeting at 5PM');
