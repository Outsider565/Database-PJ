import requests
import time
auth = ""


def get_auth_admin(user_id='test', password='passwd'):
    resp = requests.post("http://localhost/auth/admin", headers={'Content-type': 'application/json'},
                         json={"user_id": user_id, "password": password})
    print(resp.content)
    return resp.json()


def get_auth_wx(code, avatarUrl="32", nickName="232",
                gender=1, city="12", province="shan", country="CN"):
    resp = requests.post("http://localhost/auth/wx", headers={'Content-type': 'application/json'}, json={"code": code,
                         "avatar_url": avatarUrl, "nick_name": nickName, "gender": gender, "city": city, "province": province, "country": country})
    print(resp.content)
    return resp.json()


def search_class_strict(max_num, course_name="", course_id="", year="", semester="", teacher_name=""):
    resp = requests.get(
        f"http://localhost/api/course_search_strict?max_num={max_num}&course_name={course_name}&course_id={course_id}&year={year}&semester={semester}&teacher_name={teacher_name}")
    return resp.content.decode("utf-8")


def search_class(max_num, course_name="", course_id="", year="", semester="", teacher_name=""):
    resp = requests.get(
        f"http://localhost/api/course_search?max_num={max_num}&course_name={course_name}&course_id={course_id}&year={year}&semester={semester}&teacher_name={teacher_name}", headers={"token": auth})
    return resp.content.decode("utf-8")


def empty_classroom(building, date, index_str):
    resp = requests.get(
        f"http://localhost/api/empty_classroom?building={building}&date={date}&index={index_str}")
    return resp.content.decode("utf-8")


def course_time(course_id):
    resp = requests.get(
        f"http://localhost/api/course_time_loc?course_id={course_id}")
    return resp.content.decode("utf-8")


def get_classes(date, classroom_no):
    resp = requests.get(
        f"http://localhost/api/get_classes?date={date}&classroom_no={classroom_no}")
    return resp.content.decode("utf-8")


def take_course(course_id, year, semester):
    resp = requests.get(
        f"http://localhost/p_api/take_course?course_id={course_id}&year={year}&semester={semester}", headers={"token": auth})
    return resp.content.decode("utf-8")


def drop_course(course_id, year, semester):
    resp = requests.get(
        f"http://localhost/p_api/drop_course?course_id={course_id}&year={year}&semester={semester}", headers={"token": auth})
    return resp.content.decode("utf-8")

def all_course(year="",semester=""):
    resp = requests.get(
        f"http://localhost/p_api/all_course?year={year}&semester={semester}", headers={"token": auth})
    return resp.content.decode("utf-8")

def start_study(classroom_no):
    resp = requests.get(
        f"http://localhost/p_api/start_study?classroom_no={classroom_no}", headers={"token": auth})
    return resp.content.decode("utf-8")


def end_study():
    resp = requests.get(
        f"http://localhost/p_api/end_study", headers={"token": auth})
    return resp.content.decode("utf-8")


def study_time():
    resp = requests.get(
        f"http://localhost/p_api/study_time", headers={"token": auth})
    return resp.content.decode("utf-8")


def test_empty_classroom():
    print(empty_classroom("H2", date="2021-06-09", index_str="1"))
    print(empty_classroom("H2", date="2021-06-09", index_str="2"))
    print(empty_classroom("H2", date="2021-06-09", index_str="3"))
    print(empty_classroom("H2", date="2021-06-09", index_str="2"))
    print(empty_classroom("H2", date="2021-06-09", index_str="4"))
    print(empty_classroom("H2", date="2021-06-09", index_str="5"))
    print(empty_classroom("H2", date="2021-06-09", index_str="1,2,3"))
    print(empty_classroom("H2", date="2021-06-09", index_str="4,5"))


def test_search_class():
    print(search_class(max_num=10, course_name="数学分析"))
    print(search_class(max_num=10, course_id="MATH120017"))
    print(search_class_strict(max_num=10, course_id="MATH120017"))
    print(search_class_strict(max_num=10, course_id="MATH120017.04"))


def test_get_course_time():
    print(course_time("MATH120017.04"))


def test_get_classes():
    print(get_classes("2021-06-08", "H3109"))


def test_take_course():
    print(take_course("MATH120017.04", 2021, "Spring"))
    print(take_course("MATH120017.04", 2021, "Spring"))
    print(take_course("COMP130136.01", 2021, "Spring"))


def test_drop_course():
    print(drop_course("MATH120017.04", 2021, "Spring"))
    print(drop_course("MATH120017.04", 2021, "Spring"))
    print(drop_course("COMP130136.01", 2021, "Spring"))


def test_study():
    print(start_study("H3109"))
    print(start_study("H3108"))
    time.sleep(5)
    print(end_study())
    print(end_study())
    print(study_time())

def test_all_course():
    print(take_course("MATH120017.04", 2021, "Spring"))
    print(all_course())
    print(take_course("COMP130136.01", 2021, "Spring"))
    print(all_course(2021,"Spring"))
    print(all_course())
    print(drop_course("MATH120017.04", 2021, "Spring"))
    print(all_course())
    print(drop_course("COMP130136.01", 2021, "Spring"))
    print(all_course())



if __name__ == "__main__":
    auth = bytes(get_auth_wx("test")['token'], encoding="utf-8")
    print(auth)
    # test_get_course_time()
    # print(search_class(5,course_name="计算机"))
    # test_take_course()
    # test_drop_course()
    # test_study()
    # test_get_classes)
    test_all_course()
