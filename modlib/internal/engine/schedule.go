package engine

import "C"
import (
	"fmt"
	"unsafe"
)

const (
	taskSharedInvalid = iota
	taskSharedWait
	taskSharedWaitFaceEnemy
	taskSharedWaitPvs
	taskSharedSuggestState
	taskSharedWalkToTarget
	taskSharedRunToTarget
	taskSharedMoveToTargetRange
	taskSharedGetPathToEnemy
	taskSharedGetPathToEnemyLkp
	taskSharedGetPathToEnemyCorpse
	taskSharedGetPathToLeader
	taskSharedGetPathToSpot
	taskSharedGetPathToTarget
	taskSharedGetPathToHintnode
	taskSharedGetPathToLastposition
	taskSharedGetPathToBestsound
	taskSharedGetPathToBestscent
	taskSharedRunPath
	taskSharedWalkPath
	taskSharedStrafePath
	taskSharedClearMoveWait
	taskSharedStoreLastposition
	taskSharedClearLastposition
	taskSharedPlayActiveIdle
	taskSharedFindHintnode
	taskSharedClearHintnode
	taskSharedSmallFlinch
	taskSharedFaceIdeal
	taskSharedFaceRoute
	taskSharedFaceEnemy
	taskSharedFaceHintnode
	taskSharedFaceTarget
	taskSharedFaceLastposition
	taskSharedRangeAttack1
	taskSharedRangeAttack2
	taskSharedMeleeAttack1
	taskSharedMeleeAttack2
	taskSharedReload
	taskSharedRangeAttack1Noturn
	taskSharedRangeAttack2Noturn
	taskSharedMeleeAttack1Noturn
	taskSharedMeleeAttack2Noturn
	taskSharedReloadNoturn
	taskSharedSpecialAttack1
	taskSharedSpecialAttack2
	taskSharedCrouch
	taskSharedStand
	taskSharedGuard
	taskSharedStepLeft
	taskSharedStepRight
	taskSharedStepForward
	taskSharedStepBack
	taskSharedDodgeLeft
	taskSharedDodgeRight
	taskSharedSoundAngry
	taskSharedSoundDeath
	taskSharedSetActivity
	taskSharedSetSchedule
	taskSharedSetFailSchedule
	taskSharedClearFailSchedule
	taskSharedPlaySequence
	taskSharedPlaySequenceFaceEnemy
	taskSharedPlaySequenceFaceTarget
	taskSharedSoundIdle
	taskSharedSoundWake
	taskSharedSoundPain
	taskSharedSoundDie
	taskSharedFindCoverFromBestSound
	taskSharedFindCoverFromEnemy
	taskSharedFindLateralCoverFromEnemy
	taskSharedFindNodeCoverFromEnemy
	taskSharedFindNearNodeCoverFromEnemy
	taskSharedFindFarNodeCoverFromEnemy
	taskSharedFindCoverFromOrigin
	taskSharedEat
	taskSharedDie
	taskSharedWaitForScript
	taskSharedPlayScript
	taskSharedEnableScript
	taskSharedPlantOnScript
	taskSharedFaceScript
	taskSharedWaitRandom
	taskSharedWaitIndefinite
	taskSharedStopMoving
	taskSharedTurnLeft
	taskSharedTurnRight
	taskSharedRemember
	taskSharedForget
	taskSharedWaitForMovement
	lastCommonTask
)

var taskNameByID = map[int]string{
	taskSharedInvalid:                    "Invalid",
	taskSharedWait:                       "Wait",
	taskSharedWaitFaceEnemy:              "WaitFaceEnemy",
	taskSharedWaitPvs:                    "WaitPvs",
	taskSharedSuggestState:               "SuggestState",
	taskSharedWalkToTarget:               "WalkToTarget",
	taskSharedRunToTarget:                "RunToTarget",
	taskSharedMoveToTargetRange:          "MoveToTargetRange",
	taskSharedGetPathToEnemy:             "GetPathToEnemy",
	taskSharedGetPathToEnemyLkp:          "GetPathToEnemyLkp",
	taskSharedGetPathToEnemyCorpse:       "GetPathToEnemyCorpse",
	taskSharedGetPathToLeader:            "GetPathToLeader",
	taskSharedGetPathToSpot:              "GetPathToSpot",
	taskSharedGetPathToTarget:            "GetPathToTarget",
	taskSharedGetPathToHintnode:          "GetPathToHintnode",
	taskSharedGetPathToLastposition:      "GetPathToLastposition",
	taskSharedGetPathToBestsound:         "GetPathToBestsound",
	taskSharedGetPathToBestscent:         "GetPathToBestscent",
	taskSharedRunPath:                    "RunPath",
	taskSharedWalkPath:                   "WalkPath",
	taskSharedStrafePath:                 "StrafePath",
	taskSharedClearMoveWait:              "ClearMoveWait",
	taskSharedStoreLastposition:          "StoreLastposition",
	taskSharedClearLastposition:          "ClearLastposition",
	taskSharedPlayActiveIdle:             "PlayActiveIdle",
	taskSharedFindHintnode:               "FindHintnode",
	taskSharedClearHintnode:              "ClearHintnode",
	taskSharedSmallFlinch:                "SmallFlinch",
	taskSharedFaceIdeal:                  "FaceIdeal",
	taskSharedFaceRoute:                  "FaceRoute",
	taskSharedFaceEnemy:                  "FaceEnemy",
	taskSharedFaceHintnode:               "FaceHintnode",
	taskSharedFaceTarget:                 "FaceTarget",
	taskSharedFaceLastposition:           "FaceLastposition",
	taskSharedRangeAttack1:               "RangeAttack1",
	taskSharedRangeAttack2:               "RangeAttack2",
	taskSharedMeleeAttack1:               "MeleeAttack1",
	taskSharedMeleeAttack2:               "MeleeAttack2",
	taskSharedReload:                     "Reload",
	taskSharedRangeAttack1Noturn:         "RangeAttack1Noturn",
	taskSharedRangeAttack2Noturn:         "RangeAttack2Noturn",
	taskSharedMeleeAttack1Noturn:         "MeleeAttack1Noturn",
	taskSharedMeleeAttack2Noturn:         "MeleeAttack2Noturn",
	taskSharedReloadNoturn:               "ReloadNoturn",
	taskSharedSpecialAttack1:             "SpecialAttack1",
	taskSharedSpecialAttack2:             "SpecialAttack2",
	taskSharedCrouch:                     "Crouch",
	taskSharedStand:                      "Stand",
	taskSharedGuard:                      "Guard",
	taskSharedStepLeft:                   "StepLeft",
	taskSharedStepRight:                  "StepRight",
	taskSharedStepForward:                "StepForward",
	taskSharedStepBack:                   "StepBack",
	taskSharedDodgeLeft:                  "DodgeLeft",
	taskSharedDodgeRight:                 "DodgeRight",
	taskSharedSoundAngry:                 "SoundAngry",
	taskSharedSoundDeath:                 "SoundDeath",
	taskSharedSetActivity:                "SetActivity",
	taskSharedSetSchedule:                "SetSchedule",
	taskSharedSetFailSchedule:            "SetFailSchedule",
	taskSharedClearFailSchedule:          "ClearFailSchedule",
	taskSharedPlaySequence:               "PlaySequence",
	taskSharedPlaySequenceFaceEnemy:      "PlaySequenceFaceEnemy",
	taskSharedPlaySequenceFaceTarget:     "PlaySequenceFaceTarget",
	taskSharedSoundIdle:                  "SoundIdle",
	taskSharedSoundWake:                  "SoundWake",
	taskSharedSoundPain:                  "SoundPain",
	taskSharedSoundDie:                   "SoundDie",
	taskSharedFindCoverFromBestSound:     "FindCoverFromBestSound",
	taskSharedFindCoverFromEnemy:         "FindCoverFromEnemy",
	taskSharedFindLateralCoverFromEnemy:  "FindLateralCoverFromEnemy",
	taskSharedFindNodeCoverFromEnemy:     "FindNodeCoverFromEnemy",
	taskSharedFindNearNodeCoverFromEnemy: "FindNearNodeCoverFromEnemy",
	taskSharedFindFarNodeCoverFromEnemy:  "FindFarNodeCoverFromEnemy",
	taskSharedFindCoverFromOrigin:        "FindCoverFromOrigin",
	taskSharedEat:                        "Eat",
	taskSharedDie:                        "Die",
	taskSharedWaitForScript:              "WaitForScript",
	taskSharedPlayScript:                 "PlayScript",
	taskSharedEnableScript:               "EnableScript",
	taskSharedPlantOnScript:              "PlantOnScript",
	taskSharedFaceScript:                 "FaceScript",
	taskSharedWaitRandom:                 "WaitRandom",
	taskSharedWaitIndefinite:             "WaitIndefinite",
	taskSharedStopMoving:                 "StopMoving",
	taskSharedTurnLeft:                   "TurnLeft",
	taskSharedTurnRight:                  "TurnRight",
	taskSharedRemember:                   "Remember",
	taskSharedForget:                     "Forget",
	taskSharedWaitForMovement:            "WaitForMovement",
}

type Task struct {
	ID   int
	Data float32
}

// Name returns a human readable description, if this is a common task
func (task Task) Name() string {
	if 0 <= task.ID && task.ID < lastCommonTask {
		return taskNameByID[task.ID]
	}
	return fmt.Sprintf("#%v", task.ID)
}

// Schedule represents Schedule_t
type Schedule struct {
	address uintptr
}

// MakeSchedule creates Schedule from an address
func MakeSchedule(address uintptr) Schedule {
	return Schedule{address: address}
}

// Task returns Schedule_t::pTaskList[index]
func (schedule Schedule) Task(index int) *Task {
	if index < 0 || index >= schedule.TaskCount() {
		return nil
	}
	base := *(*uintptr)(unsafe.Pointer(schedule.address))
	// sizeof(Task_t) == 8
	elem := base + uintptr(8*index)
	return &Task{ID: int(*(*int32)(unsafe.Pointer(elem))), Data: *(*float32)(unsafe.Pointer(elem + 4))}
}

// TaskCount returns Schedule_t::cTasks
func (schedule Schedule) TaskCount() int {
	return int(*(*int32)(unsafe.Pointer(schedule.address + 4)))
}

// Name returns Schedule_t::pName
func (schedule Schedule) Name() string {
	return C.GoString(*(**C.char)(unsafe.Pointer(schedule.address + 16)))
}
