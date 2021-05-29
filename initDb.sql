create table GeneralInfo(
    base_id int,
    last_user_id int,
    last_group_id int,
    last_doc_id int
);

insert into GeneralInfo values (0, 0, 0, 0);

create table Users(
    id int primary key,
    login text
);

create table Password(
    id int primary key,
    password bytea,
    constraint fr_user_id foreign key(id) references Users(id)
);

create table Docs(
    id int primary key,
    creator_id int,
    text text,
    public_access_type int,
    constraint fr_creator_id foreign key(creator_id) references Users(id)
);

create table Groups1(
    id int primary key,
    creator_id int,
    name text,
    constraint fr_creator_id foreign key(creator_id) references Users(id)
);


create table DocGroupRestriction(
    doc_id int,
    group_id int,
    type int,
    constraint fr_doc_id foreign key(doc_id) references Docs(id),
    constraint fr_group_id foreign key(group_id) references Groups1(id),
    primary key(doc_id, group_id)
);

create table DocMemberRestriction(
    doc_id int,
    member_id int,
    type int,
    constraint fr_doc_id foreign key(doc_id) references Docs(id),
    constraint fr_member_id foreign key(member_id) references Users(id),
    primary key(doc_id, member_id)
);

create table GroupMember(
    group_id int,
    member_id int,
    constraint fr_group_id foreign key(group_id) references Groups1(id),
    constraint fr_member_id foreign key(member_id) references Users(id),
    primary key(group_id, member_id)
);