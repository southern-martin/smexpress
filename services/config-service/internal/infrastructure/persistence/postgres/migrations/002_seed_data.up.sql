SET search_path TO imcs_config;

INSERT INTO country_configs (country_code, country_name, currency_code, currency_symbol, timezone, date_format, weight_unit, dimension_unit, locale) VALUES
('AU', 'Australia',      'AUD', '$',    'Australia/Sydney',      'DD/MM/YYYY', 'kg', 'cm', 'en'),
('CA', 'Canada',         'CAD', '$',    'America/Toronto',       'MM/DD/YYYY', 'lb', 'in', 'en'),
('DE', 'Germany',        'EUR', '€',    'Europe/Berlin',         'DD.MM.YYYY', 'kg', 'cm', 'de'),
('FR', 'France',         'EUR', '€',    'Europe/Paris',          'DD/MM/YYYY', 'kg', 'cm', 'fr'),
('HK', 'Hong Kong',      'HKD', '$',    'Asia/Hong_Kong',        'DD/MM/YYYY', 'kg', 'cm', 'en'),
('IN', 'India',          'INR', '₹',    'Asia/Kolkata',          'DD/MM/YYYY', 'kg', 'cm', 'en'),
('KR', 'South Korea',    'KRW', '₩',    'Asia/Seoul',            'YYYY-MM-DD', 'kg', 'cm', 'ko'),
('MA', 'Morocco',        'MAD', 'د.م.', 'Africa/Casablanca',     'DD/MM/YYYY', 'kg', 'cm', 'fr'),
('NL', 'Netherlands',    'EUR', '€',    'Europe/Amsterdam',      'DD-MM-YYYY', 'kg', 'cm', 'nl'),
('NZ', 'New Zealand',    'NZD', '$',    'Pacific/Auckland',      'DD/MM/YYYY', 'kg', 'cm', 'en'),
('UK', 'United Kingdom', 'GBP', '£',    'Europe/London',         'DD/MM/YYYY', 'kg', 'cm', 'en'),
('VN', 'Vietnam',        'VND', '₫',    'Asia/Ho_Chi_Minh',      'DD/MM/YYYY', 'kg', 'cm', 'vi'),
('ZA', 'South Africa',   'ZAR', 'R',    'Africa/Johannesburg',   'DD/MM/YYYY', 'kg', 'cm', 'en');

INSERT INTO sequences (country_code, sequence_type, prefix, current_value, format_pattern) VALUES
('AU', 'shipment',    'AU', 0, '{prefix}{value}'),
('AU', 'invoice',     'INV-AU', 0, '{prefix}{value}'),
('AU', 'credit_note', 'CN-AU', 0, '{prefix}{value}'),
('UK', 'shipment',    'UK', 0, '{prefix}{value}'),
('UK', 'invoice',     'INV-UK', 0, '{prefix}{value}'),
('UK', 'credit_note', 'CN-UK', 0, '{prefix}{value}');
