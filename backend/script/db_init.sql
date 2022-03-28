create or replace view take_course as
    select *
    from users natural join takes
    order by takes.year;

create or replace function student_num_add1(course_id_ TEXT,year_ BIGINT,semester_ text) returns void
as $studentnumadd1$
begin
    update courses as c
    set student_num=student_num+1
    where c.course_id=course_id_ and c.year=year_ and c.semester = semester_;
    return;
end;
$studentnumadd1$language plpgsql;

create or replace function student_num_minus1(course_id_ TEXT,year_ BIGINT,semester_ text) returns void
as $studentnumadd1$
begin
    update courses as c
    set student_num=student_num-1
    where c.course_id=course_id_ and c.year=year_ and c.semester = semester_;
    return;
end;
$studentnumadd1$language plpgsql;

-- 在postgresql11中，我们课本上所讲的procedure和function是同义词，因此这里只是为了题目要求写出procedure，只需要把function稍加替换即可
create or replace procedure student_num_add1_p(course_id_ TEXT,year_ BIGINT,semester_ text)
as $studentnumadd1$
begin
    update courses as c
    set student_num=student_num+1
    where c.course_id=course_id_ and c.year=year_ and c.semester = semester_;
    return;
end;
$studentnumadd1$language plpgsql;

create or replace procedure student_num_minus1_p(course_id_ TEXT,year_ BIGINT,semester_ text)
as $studentnumadd1$
begin
    update courses as c
    set student_num=student_num-1
    where c.course_id=course_id_ and c.year=year_ and c.semester = semester_;
    return;
end;
$studentnumadd1$language plpgsql;

-- 选课就加人
drop trigger if exists take_course_trigger;
create trigger take_course_trigger
    after insert on takes
    referencing new row as nrow
    for each row
    execute procedure student_num_add1(nrow.course_id,nrow.year,nrow.semester);
-- 退课就减人
drop trigger if exists drop_course_trigger;
create trigger drop_course_trigger
    after delete on takes
    referencing old row as orow
    for each row
    execute procedure student_num_add1(orow.course_id,orow.year,orow.semester);
