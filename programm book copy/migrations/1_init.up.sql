CREATE TABLE IF NOT EXISTS users
	(
		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    	name TEXT NOT NULL,
    	email TEXT NOT NULL,
		password TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS books
	(
		bid serial PRIMARY KEY,
    	title TEXT NOT NULL,
    	author TEXT NOT NULL,
		id UUID NOT NULL
	);