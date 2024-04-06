package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

var TalentTreeSizes = [3]int{19, 19, 20}

// Start looking to refresh 5 minute totems at 4:55.
const TotemRefreshTime5M = time.Second * 295

const (
	SpellFlagShock     = core.SpellFlagAgentReserved1
	SpellFlagElectric  = core.SpellFlagAgentReserved2
	SpellFlagTotem     = core.SpellFlagAgentReserved3
	SpellFlagFocusable = core.SpellFlagAgentReserved4
)

func NewShaman(character *core.Character, talents string, totems *proto.ShamanTotems, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	shaman := &Shaman{
		Character:           *character,
		Talents:             &proto.ShamanTalents{},
		Totems:              totems,
		SelfBuffs:           selfBuffs,
		thunderstormInRange: thunderstormRange,
	}
	// shaman.waterShieldManaMetrics = shaman.NewManaMetrics(core.ActionID{SpellID: 57960})

	// core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents, TalentTreeSizes)
	// shaman.EnableManaBar()

	// // Add Shaman stat dependencies
	// shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	// shaman.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	// shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
	// shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	// // Set proper Melee Haste scaling
	// shaman.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	// if selfBuffs.Shield == proto.ShamanShield_WaterShield {
	// 	shaman.AddStat(stats.MP5, 100)
	// }

	// // When using the tier bonus for snapshotting we do not use the bonus spell
	// if totems.EnhTierTenBonus {
	// 	totems.BonusSpellpower = 0
	// }

	// //shaman.FireElemental = shaman.NewFireElemental(float64(totems.BonusSpellpower))
	return shaman
}

// Which buffs this shaman is using.
type SelfBuffs struct {
	Shield  proto.ShamanShield
	ImbueMH proto.ShamanImbue
	ImbueOH proto.ShamanImbue
}

// Indexes into NextTotemDrops for self buffs
const (
	AirTotem int = iota
	EarthTotem
	FireTotem
	WaterTotem
)

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	thunderstormInRange bool // flag if thunderstorm will be in range.

	Talents   *proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems *proto.ShamanTotems

	// The expiration time of each totem (earth, air, fire, water).
	TotemExpirations [4]time.Duration

	LightningBolt         *core.Spell
	LightningBoltOverload *core.Spell

	ChainLightning          *core.Spell
	ChainLightningHits      []*core.Spell
	ChainLightningOverloads []*core.Spell

	LavaBurst         *core.Spell
	LavaBurstOverload *core.Spell
	FireNova          *core.Spell
	LavaLash          *core.Spell
	Stormstrike       *core.Spell

	LightningShield     *core.Spell
	LightningShieldAura *core.Aura

	Earthquake   *core.Spell
	Thunderstorm *core.Spell

	EarthShock *core.Spell
	FlameShock *core.Spell
	FrostShock *core.Spell

	FeralSpirit *core.Spell
	//SpiritWolves *SpiritWolves

	//FireElemental      *FireElemental
	FireElementalTotem *core.Spell

	MagmaTotem           *core.Spell
	ManaSpringTotem      *core.Spell
	HealingStreamTotem   *core.Spell
	SearingTotem         *core.Spell
	StrengthOfEarthTotem *core.Spell
	//TotemicWrath         *core.Spell
	TremorTotem      *core.Spell
	StoneskinTotem   *core.Spell
	WindfuryTotem    *core.Spell
	WrathOfAirTotem  *core.Spell
	FlametongueTotem *core.Spell

	UnleashLife  *core.Spell
	UnleashFlame *core.Spell
	UnleashFrost *core.Spell
	UnleashWind  *core.Spell

	MaelstromWeaponAura *core.Aura

	// Healing Spells
	tidalWaveProc          *core.Aura
	ancestralHealingAmount float64
	AncestralAwakening     *core.Spell
	LesserHealingWave      *core.Spell
	HealingWave            *core.Spell
	ChainHeal              *core.Spell
	Riptide                *core.Spell
	EarthShield            *core.Spell

	waterShieldManaMetrics *core.ResourceMetrics

	hasHeroicPresence bool
}

// Implemented by each Shaman spec.
type ShamanAgent interface {
	core.Agent

	// The Shaman controlled by this Agent.
	GetShaman() *Shaman
}

func (shaman *Shaman) GetCharacter() *core.Character {
	return &shaman.Character
}

func (shaman *Shaman) HasPrimeGlyph(glyph proto.ShamanPrimeGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}
func (shaman *Shaman) HasMajorGlyph(glyph proto.ShamanMajorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}
func (shaman *Shaman) HasMinorGlyph(glyph proto.ShamanMinorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}

func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {

	if shaman.Totems.Fire != proto.FireTotem_NoFireTotem {
		raidBuffs.TotemicWrath = true
	}

	if shaman.Totems.Fire == proto.FireTotem_FlametongueTotem {
		raidBuffs.FlametongueTotem = true
	}

	if shaman.Totems.Water == proto.WaterTotem_ManaSpringTotem {
		raidBuffs.ManaSpringTotem = true
	}

	switch shaman.Totems.Air {
	case proto.AirTotem_WrathOfAirTotem:
		raidBuffs.WrathOfAirTotem = true
	case proto.AirTotem_WindfuryTotem:
		raidBuffs.WindfuryTotem = true
	}

	switch shaman.Totems.Earth {
	case proto.EarthTotem_StrengthOfEarthTotem:
		raidBuffs.StrengthOfEarthTotem = true
	case proto.EarthTotem_StoneskinTotem:
		raidBuffs.StoneskinTotem = true
	}

	// if shaman.Talents.UnleashedRage > 0 {
	// 	raidBuffs.UnleashedRage = true
	// }

	if shaman.Talents.ElementalOath == 1 {
		raidBuffs.ElementalOath = proto.TristateEffect_TristateEffectRegular
	} else if shaman.Talents.ElementalOath == 2 {
		raidBuffs.ElementalOath = proto.TristateEffect_TristateEffectImproved
	}
}
func (shaman *Shaman) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if shaman.Talents.ManaTideTotem {
		partyBuffs.ManaTideTotems++
	}

	shaman.hasHeroicPresence = partyBuffs.HeroicPresence
}

func (shaman *Shaman) Initialize() {
	// shaman.registerChainLightningSpell()
	// shaman.registerFeralSpirit()
	// shaman.registerFireElementalTotem()
	// shaman.registerFireNovaSpell()
	shaman.registerLavaBurstSpell()
	// shaman.registerLavaLashSpell()
	shaman.registerLightningBoltSpell()
	// shaman.registerLightningShieldSpell()
	// shaman.registerMagmaTotemSpell()
	//shaman.registerManaSpringTotemSpell()
	//shaman.registerHealingStreamTotemSpell()
	// shaman.registerSearingTotemSpell()
	shaman.registerShocks()
	// shaman.registerStormstrikeSpell()
	//shaman.registerStrengthOfEarthTotemSpell()
	shaman.registerThunderstormSpell()
	// shaman.registerTotemicWrathSpell()
	//shaman.registerFlametongueTotemSpell()
	//shaman.registerTremorTotemSpell()
	//shaman.registerStoneskinTotemSpell()
	//shaman.registerWindfuryTotemSpell()
	//shaman.registerWrathOfAirTotemSpell()
	shaman.registerUnleashElements()
	shaman.registerEarthquakeSpell()

	// // This registration must come after all the totems are registered
	//shaman.registerCallOfTheElements()

	shaman.registerBloodlustCD()

	// shaman.NewTemporaryStatsAura("DC Pre-Pull SP Proc", core.ActionID{SpellID: 60494}, stats.Stats{stats.SpellPower: 765}, time.Second*10)
}

func (shaman *Shaman) RegisterHealingSpells() {
	// shaman.registerAncestralHealingSpell()
	// shaman.registerLesserHealingWaveSpell()
	// shaman.registerHealingWaveSpell()
	// shaman.registerRiptideSpell()
	// shaman.registerEarthShieldSpell()
	// shaman.registerChainHealSpell()

	// if shaman.Talents.TidalWaves > 0 {
	// 	shaman.tidalWaveProc = shaman.GetOrRegisterAura(core.Aura{
	// 		Label:    "Tidal Wave Proc",
	// 		ActionID: core.ActionID{SpellID: 53390},
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Deactivate(sim)
	// 		},
	// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
	// 			shaman.HealingWave.CastTimeMultiplier *= 0.7
	// 			shaman.LesserHealingWave.BonusCritRating += core.CritRatingPerCritChance * 25
	// 		},
	// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 			shaman.HealingWave.CastTimeMultiplier /= 0.7
	// 			shaman.LesserHealingWave.BonusCritRating -= core.CritRatingPerCritChance * 25
	// 		},
	// 		MaxStacks: 2,
	// 	})
	// }
}

func (shaman *Shaman) Reset(sim *core.Simulation) {

}

func (shaman *Shaman) ElementalFuryCritMultiplier(secondary float64) float64 {
	elementalBonus := 0.0

	if shaman.Spec == proto.Spec_SpecElementalShaman {
		elementalBonus = 1.0
	}

	elementalBonus += secondary

	return shaman.SpellCritMultiplier(1, elementalBonus)
}

func (shaman *Shaman) GetOverloadChance() float64 {
	overloadChance := 0.0

	if shaman.Spec == proto.Spec_SpecElementalShaman {
		//  TODO: Get mastery bonus
		masteryBonus := 0.0
		overloadChance = 0.16 + masteryBonus
	}

	return overloadChance
}
