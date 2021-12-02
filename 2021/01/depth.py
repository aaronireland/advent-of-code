import sys


def read_depth_file(file_path, windowed=False):
    depths = []
    windows = {0: [], 1: [], 2: []}
    row = 1
    with open(file_path, 'r') as input:
        for depth in input.readlines():
            depth = int(depth)

            if not windowed:
                depths.append(depth)
                continue

            windows[row % 3].append(depth)
            if row - 1 > 0:
                windows[(row - 1) % 3].append(depth)
            if row - 2 > 0:
                windows[(row - 2) % 3].append(depth)

            print(f"row {row}: depth={depth}\t\t{windows}")
            for window, readings in windows.items():
                if len(readings) > 2:
                    depths.append(sum(d for d in readings))
                    windows[window] = []

            row += 1

    return depths


def increases_count(depths):
    previous = None
    topology = {'increases': 0, 'decreases': 0, 'no change': 0}
    # print(f"{len(depths)} depth readings: {depths}")
    print(f"{len(depths)} depth readings")
    for depth in depths:
        # print(f"Current depth: {depth}, previous depth: {previous}")
        if previous is not None:
            change = depth - previous
            if change > 0:
                # print("increased!")
                topology['increases'] += 1
            if change < 0:
                # print("decreased!")
                topology['decreases'] += 1
            if change == 0:
                topology['no change'] += 1

        previous = depth
    print(f"increases: {topology['increases']}, decreases: {topology['decreases']}, no change: {topology['no change']}")
    return topology['increases']


if __name__ == "__main__":
    windowed = False
    file_path = "sonar.txt"
    for i, arg in enumerate(sys.argv):
        if i > 0:
            if arg == "windowed":
                windowed = True
            else:
                file_path = arg

    depths = read_depth_file(file_path, windowed)
    answer = increases_count(depths)
    print(answer)
