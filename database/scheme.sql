CREATE TABLE links (
	link char(10) CONSTRAINT link_pk PRIMARY KEY,
	original_url varchar(2048) NOT NULL,
	
	CONSTRAINT original_url_unique UNIQUE (original_url)
);