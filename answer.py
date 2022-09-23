# main idea:
# because the max digit is 6, we can make simpler solution
# the idea is we create 3 base parse function (parse 1 digit to 3 digits)
# for 4 digits number, parse first digit, add "ribu" then parse last 3 digits
# for 5 digits number, parse first two digits, add "ribu" then parse last 3 digits
# for 6 digits number, parse first three digits, add "ribu" then parse last 3 digits
# note: in Bahasa, digit 0 will not be pronounced except for zero, so we should ignore it and just parse the next digits

single_digit = [
    "nol", "satu", "dua", "tiga", "empat", "lima", "enam", "tujuh", "delapan", "sembilan"
]

# translate one digit number to Bahasa
def parse_satuan(number):
    return single_digit[int(number)]


# translate number 11 -- 19 to Bahasa
def parse_belasan(number):
    answer = ""
    # special case
    if number == "10":
        answer += "sepuluh"
    elif number == "11":
        answer += "sebelas"
    else:
        answer += parse_satuan(number[1]) + " belas"
    return answer

# transalte number 10 -- 99
def parse_puluhan(number):
    answer = ""
    if number[0] == "0":
        answer += parse_satuan(number[1]) if number[1] != "0" else ""
    # because 11 -- 19 is special case, we need to make separate function
    elif number[0] == "1":
        answer += parse_belasan(number)
    else:
        answer += parse_satuan(number[0]) + " puluh"
        answer += " " + parse_satuan(number[1]) if number[1] != "0" else ""
    return answer

# transalte number 100 -- 999
def parse_ratusan(number):
    answer = ""
    if number[0] == "0":
        answer += parse_puluhan(number[1:])
        return answer
    elif number[0] == "1":
        answer += "seratus"
    else:
        answer += parse_satuan(number[0]) + " ratus"
    
    answer += " " + parse_puluhan(number[1:]) if parse_puluhan(number[1:]) != "" else ""
    return answer

# transalte number 1000 -- 9999
def parse_ribuan(number):
    answer = ""
    if number[0] == "1":
        answer += "seribu"
    elif number[0] != "0":
        answer += parse_satuan(number[0]) + " ribu"
    answer += " " + parse_ratusan(number[1:]) if parse_ratusan(number[1:]) != "" else ""
    return answer

# transalte number 10000 -- 99999
def parse_puluh_ribuan(number):
    answer = ""
    if number[0] == "0":
        answer += parse_ribuan(number[1:])
    else:
        answer += parse_puluhan(number[:2]) + " ribu"
    answer += " " + parse_ratusan(number[2:]) if parse_ratusan(number[2:]) != "" else ""
    return answer

# transalte number 100000 -- 999999
def parse_ratus_ribuan(number):
    answer = ""
    if number[0] == "0":
        answer += parse_puluh_ribuan(number[1:])
    else:
        answer += parse_ratusan(number[:3]) + " ribu"
    answer += " " + parse_ratusan(number[3:]) if parse_ratusan(number[3:]) != "" else ""
    return answer

def convert(number):
    # remove leading zeros
    while number[0] == "0" and len(number) != 1:
        number = number[1:]
    
    length = len(number)
    if length == 1:
        return parse_satuan(number)
    elif length == 2:
        return parse_puluhan(number)
    elif length == 3:
        return parse_ratusan(number)
    elif length == 4:
        return parse_ribuan(number)
    elif length == 5:
        return parse_puluh_ribuan(number)
    elif length == 6:
        return parse_ratus_ribuan(number)
    else:
        print("this program does not support integer more than 6 digits")
        return ""

n = input("Masukkan angka maksimum 6 digit: ")
if n == "":
    print("input cant be null")

converted = convert(n)
print(converted)
