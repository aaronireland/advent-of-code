import sys

from collections import defaultdict


def points(seg):
    p1, p2 = (seg[0], seg[1],) if ((seg[0][0] < seg[1][0]) or (seg[0][0] == seg[1][0] and seg[0][1] <= seg[1][1])) else (seg[1], seg[0],)
    dx, dy = p2[0] - p1[0], p2[1] - p1[1]

    # DIAGONAL
    if dx != 0 and dy != 0:
        m = dx / dy
        if abs(int(m)) != 1:
            raise Exception("diagonal segments must have slope of 1, m={m}")

        x0 = p1[0]
        y0 = p1[1]
        for x in range(dx + 1):
            yield(x0 + x, y0 + int(x * m))

    # HORIZONTAL
    if dx != 0 and dy == 0:
        y = seg[0][1]
        x0 = min(seg[0][0], seg[1][0])
        for x in range(dx + 1):
            yield (x0 + x, y,)

    # VERTICAL
    if dy != 0 and dx == 0:
        x = seg[0][0]
        y0 = min(seg[0][1], seg[1][1])
        for y in range(dy + 1):
            yield (x, y0 + y)

    # SINGLE POINT
    if dx == 0 and dy == 0:
        yield p1


def filter_horiz_vert_only(lines):
    hv_lines = []
    for line in lines:
        if line[0][0] == line[1][0] or line[0][1] == line[1][1]:
            hv_lines.append(line)
    return hv_lines


def read_file(file_path):
    line_segments = []
    with open(file_path, 'r') as input:
        for line in input.readlines():
            coords = line.split(' -> ')
            coord_a = coords[0].split(',')
            coord_b = coords[1].split(',')
            line_segments.append(((int(coord_a[0]), int(coord_a[1]),), (int(coord_b[0]), int(coord_b[1]),)))
    return line_segments


if __name__ == '__main__':
    file_path = 'input.txt'
    diagonals = False

    for i, arg in enumerate(sys.argv):
        if i > 0:
            if arg == "diagonals":
                diagonals = True
            else:
                file_path = arg

    lines = read_file(file_path) if diagonals else filter_horiz_vert_only(read_file(file_path))
    thermal_map = defaultdict(int)

    for line in lines:
        for coord in points(line):
            thermal_map[coord] += 1

    hotspots = sum(1 for coord, count in thermal_map.items() if count > 1)

    print(f"Thermal scan shows {hotspots} hot spots!")
