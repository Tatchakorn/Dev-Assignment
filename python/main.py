from dataclasses import dataclass
from datetime import datetime
from enum import Enum
from typing import Tuple

DATETIME_FORMAT = r'%Y-%m-%d'
START_DATE = datetime(2564,6,1)
END_DATE = datetime(2564,8,31)
SENIOR_LOWER_AGE = (65, 0)
CHILD_LOWER_AGE = (0, 6)
CHILD_UPPER_AGE = (2, 0)
TH_MONTHS = (
    'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน', 
    'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
)

class Gender(Enum):
    Male, Female = ('Male', 'Female')


    def __str__(self):
        return self.value


@dataclass
class Person:
    gender: Gender
    birthdate: datetime


    def __str__(self) -> str:
        return f'Gender: {str(self.gender)}, ' \
                f'birthdate: {self.th_date_format(self.birthdate)}'


    def calculate_age(self, set_date: datetime) -> Tuple[int, int]:
        '''
        returns age in (years, months) from birthdate relative to the set date.
        '''
        age_years = set_date.year - self.birthdate.year
        
        # has passed the set date's birthday for that year
        if (set_date.month < self.birthdate.month) or \
            (set_date.month == self.birthdate.month and set_date.day < self.birthdate.day):
            age_years -= 1
        
        age_months = ((set_date.month - self.birthdate.month) + 12)  % 12 # Prefer pos int
        return age_years, age_months


    def is_eligible(self) -> Tuple[bool, datetime, datetime]:
        '''
        returns if the person is able to apply for the service within the sevice period
        with start and end date of the period when the person can apply for the service.
        For senior citizens (65 years old or older) 
        [65, inf)
        For children (between 6 months and 2 years old) 
        [0.6, 2]
        If the person does not meet any of these criteria, they are deemed ineligible.
        '''

        age_at_start = self.calculate_age(START_DATE)
        age_at_end = self.calculate_age(END_DATE)
        dob = self.birthdate

        # old enough and not too old
        if  SENIOR_LOWER_AGE <= age_at_start or \
            CHILD_LOWER_AGE <= age_at_start  and  age_at_end <= CHILD_UPPER_AGE:
            return True, START_DATE, END_DATE        
        
        # will be 65yo old within the period
        elif (SENIOR_LOWER_AGE <= age_at_end):
            return True, datetime(dob.year+65, dob.month, dob.day), END_DATE
        
        # will be 6mo within the period
        elif CHILD_LOWER_AGE <= age_at_end < CHILD_UPPER_AGE:
            return True, datetime(dob.year, dob.month+6, dob.day) , END_DATE
        
        # will be 2yo within the period
        elif CHILD_UPPER_AGE < age_at_end and \
            CHILD_LOWER_AGE <= age_at_start <= CHILD_UPPER_AGE:
            return True, START_DATE, datetime(dob.year+2, dob.month, dob.day)
        
        # ineligible
        return False, None, None
    
    
    @staticmethod
    def th_date_format(date: datetime) -> str:
        '''
        Converts datetime in the format of "%Y-%m-%d" to Thai-styled date
        '''
        return f'{date.day} {TH_MONTHS[date.month - 1]} พ.ศ.{date.year}'


def main() -> None:
    example_dates = [
        ('2499-03-10', Gender.Female),  # 10 มีนาคม พ.ศ.2499 
        ('2500-10-08', Gender.Male),    # 8 ตุลาคม พ.ศ.2500 
        ('2562-07-01', Gender.Female),  # 1 กรกฎาคม พ.ศ.2562 
        ('2564-01-05', Gender.Female),  # 5 มกราคม พ.ศ.2564
    ]
    str_to_date = lambda date_str: datetime.strptime(date_str, DATETIME_FORMAT)
    persons =  [
        Person(birthdate=str_to_date(dob), gender=gender) for dob, gender in example_dates
    ]

    for i, p in enumerate(persons):
        eligible, start, end = p.is_eligible()
        eligible = 'Yes' if eligible else 'No'
        start = p.th_date_format(start) if start else None
        end = p.th_date_format(end) if end else None
        print(f'{i+1}: {p}')
        print(
            f'{eligible},\nService Start Date:\t{start}\n' \
                f'Service End Date:\t{end}\n{"-" * 10}')


if __name__ == '__main__':
    main()