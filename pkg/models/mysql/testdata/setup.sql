CREATE TABLE IF NOT EXISTS texts (
    urltoken BINARY(32) NOT NULL PRIMARY KEY,
    text TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    views INTEGER NOT NULL,
    viewlimit INTEGER NOT NULL
);

/*plainTextToken: FTR43TPBEWDCQ4B2HRCNXPSDBXFEAQ44QWC7QZ2P5D5NW3Y64UJA */
INSERT INTO texts (urltoken, text, created, expires, views, viewlimit) VALUES (
    0x3FA4941C5FDA1A71EA31A94312ADB2CA58480F98415B2A3F1D24F8F99F7C5C3C,
    'This is an example tempshare for testing purposes!',
    '2022-03-02 12:00:00',
    '2048-03-09 12:00:00',
    0,
    1
);

/*plainTextToken: HVN2JMTD5DVPODS632YXWVT6REYSXR26O7B3G5ZBQRD72IOBYTVA */
INSERT INTO texts (urltoken, text, created, expires, views, viewlimit) VALUES (
    0x87236F3ED11C646E80652DE80FB121F6315BB5BB7C649E83251DD088D2A61148,
    'This is an expired tempshare!',
    '2022-03-02 12:00:00',
    '2048-03-09 12:00:00',
    1,
    1
);
