
from django.db import models


class Admin(models.Model):
    admin_id = models.TextField(primary_key=True)
    password_hash = models.TextField()

    class Meta:
        managed = False
        db_table = 'admins'



class Building(models.Model):
    building_no = models.TextField(primary_key=True,verbose_name='教学楼号')
    building_name = models.TextField(blank=True, null=True,verbose_name='教学楼名')
    latitude = models.FloatField(blank=True, null=True,verbose_name='纬度')
    longitude = models.FloatField(blank=True, null=True,verbose_name='经度')

    class Meta:
        managed = False
        db_table = 'buildings'
        verbose_name = '教学楼'
        verbose_name_plural = '教学楼'

class Classroom(models.Model):
    classroom_no = models.TextField(primary_key=True,verbose_name='教室号')
    capacity = models.BigIntegerField(blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'classrooms'
        verbose_name = '教室'
        verbose_name_plural = '教室'


class Course(models.Model):
    course_id = models.TextField(primary_key=True,verbose_name='课程号')
    year = models.IntegerField(verbose_name='年份')
    semester = models.TextField(verbose_name='学期')
    course_name = models.TextField(blank=True, null=True,verbose_name='课程名')
    student_num = models.BigIntegerField(blank=True, null=True,verbose_name='学生数量')
    teacher_name = models.TextField(blank=True, null=True,verbose_name='任课教师')
    type = models.TextField(blank=True, null=True,verbose_name='类型')

    class Meta:
        managed = False
        db_table = 'courses'
        unique_together = (('course_id', 'year', 'semester'),)
        verbose_name = '课程'
        verbose_name_plural = '课程'


class LocateIn(models.Model):
    classroom_no = models.OneToOneField(
        Classroom, models.DO_NOTHING, db_column='classroom_no', primary_key=True,verbose_name='教室号')
    building_no = models.ForeignKey(
        Building, models.DO_NOTHING, db_column='building_no', blank=True, null=True,verbose_name='楼号')
    class Meta:
        managed = False
        db_table = 'locate_ins'
        verbose_name = '地址'
        verbose_name_plural = '地址'
    def __str__(self):
        return "Loc"


class User(models.Model):
    user_id = models.UUIDField(primary_key=True,verbose_name='用户ID')
    open_id = models.TextField(unique=True, blank=True, null=True,verbose_name='微信OpenID')
    avatar_url = models.TextField(blank=True, null=True,verbose_name='头像url')
    nick_name = models.TextField(blank=True, null=True,verbose_name='微信昵称')
    gender = models.SmallIntegerField(blank=True, null=True,verbose_name='性别')
    city = models.TextField(blank=True, null=True,verbose_name='城市')
    province = models.TextField(blank=True, null=True,verbose_name='省份')
    country = models.TextField(blank=True, null=True,verbose_name='国家')

    class Meta:
        managed = False
        db_table = 'users'
        verbose_name = '用户'
        verbose_name_plural = '用户'
class Study(models.Model):
    study_id = models.AutoField(primary_key=True,db_column="study_id",verbose_name="学习ID")
    user = models.ForeignKey(User, models.DO_NOTHING, blank=True, null=True,verbose_name="用户")
    start_time = models.DateTimeField(verbose_name='开始时间')
    time_len = models.FloatField(blank=True, null=True,verbose_name='自习时长')
    classroom_no = models.ForeignKey(Classroom, models.DO_NOTHING, db_column="classroom_no",verbose_name='教室号')

    class Meta:
        managed = False
        db_table = 'studies'
        verbose_name = '学习'
        verbose_name_plural = '学习'
class Take(models.Model):
    user_id = models.ForeignKey(
        User, models.DO_NOTHING, blank=True,primary_key=True,db_column="user_id",verbose_name='用户ID')
    course_id = models.TextField(verbose_name="课程号")
    year = models.IntegerField(verbose_name="年份")
    semester = models.TextField(verbose_name="学期")

    class Meta:
        managed = False
        db_table = 'takes'
        unique_together = (('user_id', 'course_id', 'year', 'semester'),)
        verbose_name = '上课'
        verbose_name_plural = '上课'

class TeachIn(models.Model):
    course = models.OneToOneField(Course, models.DO_NOTHING, primary_key=True,verbose_name="课程编号")
    year = models.IntegerField(verbose_name="年份")
    semester = models.TextField(verbose_name="学期")
    classroom_no = models.ForeignKey(
        Classroom, models.DO_NOTHING, blank=True, null=True, db_column="classroom_no",verbose_name="教室号")
    date = models.DateField(verbose_name="日期")
    class_index = models.IntegerField(verbose_name="节数")

    class Meta:
        managed = False
        db_table = 'teach_ins'
        unique_together = (('course', 'year', 'semester',
                            'classroom_no', 'date', 'class_index'),)
        verbose_name = '教学'
        verbose_name_plural = '教学'


