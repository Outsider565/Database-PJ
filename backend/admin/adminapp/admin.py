from django.contrib import admin

# Register your models here.
from django.contrib import admin
from .models import Building, Classroom, Course, LocateIn, Study, TeachIn, Take, User


class BuildingAdmin(admin.ModelAdmin):
    list_display = ("building_no", "building_name", "latitude", "longitude")

class ClassroomAdmin(admin.ModelAdmin):
    list_display = ("classroom_no", "capacity")


class CourseAdmin(admin.ModelAdmin):
    list_display = ("course_id", "year", "semester", "course_name", "student_num", "teacher_name", "type")
    list_filter = ("year","semester","type")


class LocateAdmin(admin.ModelAdmin):
    list_display = ("classroom_no","building_no")

class UserAdmin(admin.ModelAdmin):
    list_display = ("user_id","nick_name","gender","city")
    list_filter = ("gender",)

class StudyAdmin(admin.ModelAdmin):
    list_display = ("user","start_time","time_len","classroom_no")

class TakeAdmin(admin.ModelAdmin):
    list_display = ("user_id","course_id","year","semester")

class TeachInAdmin(admin.ModelAdmin):
    list_display = ("course","year","semester","classroom_no","date","class_index")

admin.site.register(Building, BuildingAdmin)
admin.site.register(Classroom, ClassroomAdmin)
admin.site.register(Course, CourseAdmin)
admin.site.register(LocateIn,LocateAdmin)
admin.site.register(Study,StudyAdmin)
admin.site.register(TeachIn,TeachInAdmin)
admin.site.register(Take,TakeAdmin)
admin.site.register(User,UserAdmin)
