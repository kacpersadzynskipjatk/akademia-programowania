package academy

import (
	"math"
)

type Student struct {
	Name       string
	Grades     []int
	Project    int
	Attendance []bool
}

// AverageGrade returns an average grade given a
// slice containing all grades received during a
// semester, rounded to the nearest integer.
func AverageGrade(grades []int) int {
	if len(grades) == 0 {
		return 0
	}
	var sum int
	for _, grade := range grades {
		sum += grade
	}
	result := float64(sum) / float64(len(grades))
	return int(math.Round(result))
}

// AttendancePercentage returns a percentage of class
// attendance, given a slice containing information
// whether a student was present (true) of absent (false).
//
// The percentage of attendance is represented as a
// floating-point number ranging from  0 to 1,
// with 2 digits of precision.
func AttendancePercentage(attendance []bool) float64 {
	if len(attendance) == 0 {
		return .0
	}
	var count int
	for _, attend := range attendance {
		if attend {
			count++
		}
	}
	result := float64(count) / float64(len(attendance))
	return math.Round(result*100) / 100
}

// FinalGrade returns a final grade achieved by a student,
// ranging from 1 to 5.
//
// The final grade is calculated as the average of a project grade
// and an average grade from the semester, with adjustments based
// on the student's attendance. The final grade is rounded
// to the nearest integer.
//
// If the student's attendance is below 80%, the final grade is
// decreased by 1. If the student's attendance is below 60%, average
// grade is 1 or project grade is 1, the final grade is 1.
func FinalGrade(s Student) int {
	avgGrade := AverageGrade(s.Grades)
	attendancePercentage := AttendancePercentage(s.Attendance)
	finalGrade := (float64(s.Project + avgGrade)) / 2.0
	//If a student have: avgGrade = 0, attendancePercentage = 1, s.Project = 5,
	//then he will get a finalGrade (5+0)/2 = 2.5,
	//I don't know if that is expected behavior because maybe he should get finalGrade = 1 due to lack of grades.
	//If it is a problem we should change the condition from avgGrade == 1 to avgGrade <= 1 in if statement.
	if attendancePercentage < 0.6 || s.Project == 1 || avgGrade == 1 {
		finalGrade = 1
	} else if attendancePercentage < 0.8 {
		finalGrade -= 1
	}
	return int(math.Round(finalGrade))
}

// GradeStudents returns a map of final grades for a given slice of
// Student structs. The key is a student's name and the value is a
// final grade.
func GradeStudents(students []Student) map[string]uint8 {
	studentsGrades := make(map[string]uint8)
	for _, student := range students {
		studentsGrades[student.Name] = uint8(FinalGrade(student))
	}
	return studentsGrades
}
