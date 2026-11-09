-- +goose Up
create table reverseindex (
    wordid serial,
    websiteid uuid,
    foreign key (wordid) references words(id) on delete cascade,
    foreign key (websiteid) references websites(id) on delete cascade
);

-- +goose Down
drop table reverseindex;
