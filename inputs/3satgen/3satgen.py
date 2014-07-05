import argparse
import random

parser = argparse.ArgumentParser(description='Generate inputs for 3SAT problem')
parser.add_argument('-a', '--arguments', type=int, help='The number of arguments', required=True)
parser.add_argument('-c', '--clauses', type=int, help='The number of arguments', required=True)
parser.add_argument('-r', '--rtrue', help='If this flag is set then the result is always true', action='store_true')


def get_non_zero_random_int(b):
    randomInteger = 0
    while randomInteger == 0:
        randomInteger = random.randint(-b, b)
    return randomInteger


if __name__ == '__main__':
    args = parser.parse_args()
    print('%d %d' % (args.clauses, args.arguments))

    random.seed()
    maxNumber = 2 ** args.arguments - 1
    selectedNumber = random.randint(0, maxNumber)

    for i in range(args.clauses):
        firstVar = get_non_zero_random_int(args.arguments)
        secondVar = get_non_zero_random_int(args.arguments)
        thirdVar = get_non_zero_random_int(args.arguments)

        firstRes = ((2 ** (abs(firstVar) - 1)) & selectedNumber) > 0
        firstRes = not firstRes if firstVar < 0 else firstRes
        secondRes = ((2 ** (abs(secondVar) - 1)) & selectedNumber) > 0
        secondRes = not secondRes if firstVar < 0 else secondRes
        thirdRes = ((2 ** (abs(thirdVar) - 1)) & selectedNumber) > 0
        thirdRes = not thirdRes if firstVar < 0 else thirdRes

        if not (firstRes or secondRes or thirdRes):
            firstVar = -firstVar
        print('%d %d %d' % (firstVar, secondVar, thirdVar))

    print('\n\n\nSelected Number :%d' % selectedNumber)