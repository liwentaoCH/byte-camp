create
database camp_base default character set utf8mb4 collate utf8mb4_unicode_ci;

DROP TABLE IF EXISTS t_member;
CREATE TABLE t_member
(
    USER_ID   BIGINT(20) NOT NULL COMMENT '用户ID',
    USER_NAME VARCHAR(32)  NOT NULL COMMENT '用户名',
    NICKNAME  VARCHAR(32)  NOT NULL COMMENT '用户昵称',
    PASSWORD  VARCHAR(128) NOT NULL COMMENT '密码',
    USER_TYPE TINYINT(4) NOT NULL COMMENT '用户类型;1-管理员，2-学生，3-教师',
    STATUS    TINYINT(4) NOT NULL COMMENT '用户状态;0-删除，1-正常',
    PRIMARY KEY (USER_ID)
) COMMENT = '用户表';


DROP TABLE IF EXISTS t_course;
CREATE TABLE t_course
(
    COURSE_ID    BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '课程ID',
    NAME         VARCHAR(32) NOT NULL COMMENT '课程名称',
    COURSE_STOCK INT(11) NOT NULL COMMENT '课程容量',
    TEACHER_ID   BIGINT(20) COMMENT '授课教师ID',
    PRIMARY KEY (COURSE_ID)
) COMMENT = '课程表';


DROP TABLE IF EXISTS t_student_course;
CREATE TABLE t_student_course
(
    STUDENT_ID BIGINT(20) NOT NULL COMMENT '学生ID',
    COURSE_ID  BIGINT(20) NOT NULL COMMENT '课程ID',
    PRIMARY KEY (STUDENT_ID, COURSE_ID)
) COMMENT = '学生课程关系表';


DROP TABLE IF EXISTS t_record;
CREATE TABLE t_record
(
    STUDENT_ID  BIGINT(20) NOT NULL COMMENT '学生ID',
    COURSE_ID   BIGINT(20) NOT NULL COMMENT '课程ID',
    CREATE_TIME timestamp NOT NULL COMMENT '创建时间',
) COMMENT = '学生课程关系表';

DROP TABLE IF EXISTS t_record;
CREATE TABLE `t_record`
(
    `id`         bigint AUTO_INCREMENT,
    `student_id` bigint(20) NOT NULL COMMENT '学生ID',
    `course_id`  bigi(20) NOT NULL COMMENT '课程ID',
    PRIMARY KEY (`id`)
)COMMENT = '差错记录表';


CREATE
unique INDEX idx_t_member_user_name ON t_member(USER_NAME)
CREATE
INDEX idx_t_course_teacher_id ON t_course(TEACHER_ID)


insert into t_member values(1000000000000000000,'JudgeAdmin','JudgeAdmin','$2a$10$D/CUcgjgqXAVd0aC8stqs.TTzgAY.AMm4i5n8rTr3IsBQFXwpy/Oe',1,1)