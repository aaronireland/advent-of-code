import copy
import sys
import re

MAX_ITERATIONS = 100


def read_file(file_path):
    bytes = []
    with open(file_path, 'r') as input:
        for line in input.readlines():
            byte = re.sub("[^0-1]", "", line)
            if len(byte) > 0:
                bytes.append(byte)
    return bytes


def frequency(bytes, pos):
    freq = {0: 0, 1: 0}
    for byte in bytes:
        try:
            bit = int(byte[pos])
            freq[bit] += 1
        except (IndexError, ValueError):
            continue
    return freq


def most_common_bit(freq):
    if freq[1] >= freq[0]:
        return '1'
    else:
        return '0'


def least_common_bit(freq):
    if freq[0] <= freq[1]:
        return '0'
    else:
        return '1'


def gamma_rate(input):
    most_common_bits = ''
    for pos in range(len(input[0])):
        most_common_bits += most_common_bit(frequency(input, pos))

    return most_common_bits


def epsilon_rate(gamma):
    return gamma.replace('1', '2').replace('0', '1').replace('2', '0')


def filtered_bytes(bytes, pos, most_common=True, verbose=False):
    filtered = []
    filter_bit = most_common_bit(frequency(bytes, pos)) if most_common else least_common_bit(frequency(bytes, pos))
    if verbose:
        print(f"filtering {len(bytes)} bytes: {'MOST' if most_common else 'LEAST'} common bit is {filter_bit}")

    for byte in bytes:
        msg = (f"checking {byte} at pos {pos} for filter {filter_bit}")
        if byte[pos] == filter_bit:
            filtered.append(byte)
            msg = (f"{msg}...added! {len(filtered)} bytes...")
        if verbose:
            print(f"{msg}")

    return filtered


def get_reading(input, most_common=True, verbose=False):
    iterations = 0
    bytes = copy.deepcopy(input)
    while True:
        iterations += 1
        if iterations > MAX_ITERATIONS:
            raise Exception("Exceeded max iterations without producing a valid O2 reading")
        for pos in range(len(input[0])):
            bytes = filtered_bytes(bytes, pos, most_common, verbose)
            if verbose:
                print(f"\n\nfilter at pos {pos} complete!\n{bytes}\n\n...")
            if len(bytes) == 1:
                return bytes[0]
            if len(bytes) == 0:
                raise Exception("filtering function emptied byte array!")
    raise Exception("filter failed")


def o2_reading(input, verbose=False):
    return get_reading(input, True, verbose)


def co2_reading(input, verbose=False):
    return get_reading(input, False, verbose)


if __name__ == '__main__':
    file_path = 'input.txt'
    verbose = False

    for i, arg in enumerate(sys.argv):
        if i > 0:
            if arg == "verbose":
                verbose = True
            else:
                file_path = arg

    bytes = read_file(file_path)

    # Part 1:
    gamma = gamma_rate(bytes)
    epsilon = epsilon_rate(gamma)
    power_consumption = int(gamma, 2) * int(epsilon, 2)
    print(f"gamma_rate={gamma}, epsilon_rate={epsilon}, power consumption is {power_consumption}")

    # Part 2
    o2 = o2_reading(bytes, verbose)
    co2 = co2_reading(bytes, verbose)
    rating = int(o2, 2) * int(co2, 2)

    print(f"o2 reading is {o2}, co2 reading is {co2}, life support rating is {rating}")
