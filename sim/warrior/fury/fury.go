package fury

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warrior"
)

func RegisterFuryWarrior() {
	core.RegisterAgentFactory(
		proto.Player_FuryWarrior{},
		proto.Spec_SpecFuryWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFuryWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FuryWarrior)
			if !ok {
				panic("Invalid spec value for Fury Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type FuryWarrior struct {
	*warrior.Warrior

	Options *proto.FuryWarrior_Options
}

func NewFuryWarrior(character *core.Character, options *proto.Player) *FuryWarrior {
	warOptions := options.GetFuryWarrior().Options

	war := &FuryWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			StanceSnapshot: warOptions.StanceSnapshot,
		}),
		Options: warOptions,
	}

	rbo := core.RageBarOptions{
		StartingRage:   warOptions.StartingRage,
		RageMultiplier: core.TernaryFloat64(war.Talents.EndlessRage, 1.25, 1),
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}
	if oh := war.GetOHWeapon(); oh != nil {
		rbo.OHSwingSpeed = oh.SwingSpeed
	}

	war.EnableRageBar(rbo)
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: war.TryHSOrCleave,
	})

	return war
}

func (war *FuryWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *FuryWarrior) Initialize() {
	war.Warrior.Initialize()

	if war.Options.UseRecklessness {
		war.RegisterRecklessnessCD()
	}

	if war.Options.UseShatteringThrow {
		war.RegisterShatteringThrowCD()
	}

	if war.PrimaryTalentTree == warrior.FuryTree {
		war.BerserkerStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
	} else if war.PrimaryTalentTree == warrior.FuryTree {
		war.BattleStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
	}
}

func (war *FuryWarrior) Reset(sim *core.Simulation) {
	if war.PrimaryTalentTree == warrior.FuryTree {
		war.Warrior.Reset(sim)
		war.BerserkerStanceAura.Activate(sim)
		war.Stance = warrior.BerserkerStance
	} else if war.PrimaryTalentTree == warrior.FuryTree {
		war.Warrior.Reset(sim)
		war.BattleStanceAura.Activate(sim)
		war.Stance = warrior.BattleStance
	}
}
