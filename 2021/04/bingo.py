import sys
import re


class Board:
    def __init__(self):
        self.rows = 5
        self.cols = 5
        self.called_numbers = []
        self.board = [[-1] * self.cols for i in range(self.rows)]

    def __str__(self):
        rows = len(self.board)  # Get the first dimension
        board_str = ''

        board_str += '------------- output -------------\n'

        for i in range(rows):
            board_str += ('|' + ', '.join(map(lambda x: '{0:5}'.format(x), self.board[i])) + '| \n')

        board_str += '----------------------------------'
        return board_str

    def __getitem__(self, key):
        if isinstance(key, tuple):
            i = key[0]
            j = key[1]
            return self.board[i][j]

    def __setitem__(self, key, value):
        if isinstance(key, tuple):
            i = key[0]
            j = key[1]
            if i >= self.rows or j >= self.cols:
                raise ValueError(f"max dim is {self.rows}x{self.cols}, {i}x{j} given.")
            self.board[i][j] = value

    def is_full(self):
        for r in range(self.rows):
            for c in range(self.cols):
                if self.board[r][c] < 0:
                    return False
        return True

    def call(self, num):
        self.called_numbers.append(num)

    def unmarked_numbers(self):
        unmarked = []
        for r in range(self.rows):
            for c in range(self.cols):
                if self.board[r][c] not in self.called_numbers:
                    unmarked.append(self.board[r][c])
        return unmarked

    def winning_row(self):
        for r in range(self.rows):
            row = [self.board[r][c] for c in range(self.cols)]
            matched_numbers = list(filter(lambda num: num in self.called_numbers, row))
            if len(matched_numbers) == self.cols:
                return True
        return False

    def winning_column(self):
        for c in range(self.cols):
            row = [self.board[r][c] for r in range(self.rows)]
            if len(list(filter(lambda num: num in self.called_numbers, row))) == self.rows:
                return True
        return False


def read_file(file_path):
    called_numbers = []
    boards = []
    row = -1
    board = Board()
    with open(file_path, 'r') as input:
        for line in input.readlines():
            if row == -1:
                called_numbers = [int(x) for x in line.split(',')]
                row += 1
                continue
            board_row = [int(x) for x in filter(lambda x: re.sub("[^0-9]", "", x), line.split(' '))]
            if len(board_row) < 5:
                continue
            if board.is_full():
                boards.append(board)
                board = Board()
            for j, num in enumerate(board_row):
                board[row % 5, j] = num
            row += 1
        boards.append(board)
    return called_numbers, boards


if __name__ == '__main__':
    file_path = 'input.txt'
    lose_on_purpose = False
    show_boards = False

    for i, arg in enumerate(sys.argv):
        if i > 0:
            if arg == "lose":
                lose_on_purpose = True
            elif arg == "show":
                show_boards = True
            else:
                file_path = arg

    called_nums, boards = read_file(file_path)
    winning_boards = []
    if show_boards:
        for i, board in enumerate(boards):
            print(f"BOARD #{i}:\n{board}\n\n")
    for number in called_nums:
        for i, board in enumerate(boards):
            if i in set([x[0] for x in winning_boards]):
                continue
            board.call(number)
            if board.winning_column() or board.winning_row():
                score = sum(board.unmarked_numbers()) * number
                if lose_on_purpose:
                    winning_boards.append((i, board, score,))
                    continue
                else:
                    print(f"Board {i+1} wins with a score of {score}!")
                    print(board)
                    sys.exit(0)
    i, losing_board, losing_score = winning_boards[-1]
    print(f"Board {i+1} won last with a score of {losing_score}!")
    print(losing_board)
