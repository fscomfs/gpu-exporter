from ctypes import *

lib = CDLL('./libtpu_util.so')


class TPU_MEMORY_INFO_STRU(Structure):
    _fields_ = (
        ("dev_id", c_int),
        ("chip_id", c_int),
        ("chip_mode", c_int),
        ("domain_bdf", c_int),
        ("status", c_int),
        ("card_index", c_int),
        ("chip_index_of_card", c_int),
        ("mem_used", c_int),
        ("mem_total", c_int),
        ("tpu_util", c_int),
        ("board_temp", c_int),
        ("chip_temp", c_int),
        ("board_power", c_int),
        ("tpu_power", c_int),
        ("fan_speed", c_int),
    )

def getAllDeviceInfo():
    result = []
    car_num = c_int()
    ret = lib.getTpuNum(byref(car_num))
    car_info = (TPU_MEMORY_INFO_STRU*car_num.value)()

    lib.getTpuMemInfo(byref(car_info),car_num);
    for carIndex in range(car_num.value):
        print(car_info[carIndex].dev_id)
        result.append({"dev_id":car_info[carIndex].dev_id
                          , "chip_id": car_info[carIndex].chip_id
                          , "chip_mode": car_info[carIndex].chip_mode
                          , "domain_bdf": car_info[carIndex].domain_bdf
                          , "status": car_info[carIndex].status
                          , "card_index": car_info[carIndex].card_index
                          , "chip_index_of_card": car_info[carIndex].chip_index_of_card
                          , "mem_used": car_info[carIndex].mem_used
                          , "mem_total": car_info[carIndex].mem_total
                          , "tpu_util": car_info[carIndex].tpu_util
                          , "board_temp": car_info[carIndex].board_temp
                          , "board_power": car_info[carIndex].board_power
                          , "tpu_power": car_info[carIndex].tpu_power
                          , "fan_speed": car_info[carIndex].fan_speed

                       })

    return result